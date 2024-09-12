package log_rotation

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	//backupTimeFormat 时间基础
	backupTimeFormat = "2006-01-02T15-04-05.000"
	//compressSuffix 压缩类型
	compressSuffix = ".gz"
	//最大数量
	defaultMaxSize = 100
)

// 为了看是否所有接口都实现
var _ io.WriteCloser = (*Logger)(nil)

type Logger struct {
	//文件名
	Filename string
	//每个文件最大大小
	MaxSize int
	//文件最长保存时间
	MaxSave int
	//文件保存最大数量
	MaxNumbers int
	//文件是否压缩
	Compress bool
	//原子操作 加锁
	mu sync.Mutex
	//文件类型，包含文件操作
	file *os.File
	//文件现在大小
	size int64
	//startMill 它用于确保某个操作只执行一次，即使该操作被多个 goroutine 同时请求。
	startMill sync.Once
	millCh    chan bool
}

var (
	//兆字节是MaxSize和字节之间的转换系数
	megabyte = 1024 * 1024
)

// write 主要是实现io.WriteCloser里边的write
// zap日志里边也是使用write进行输出
// 主要是改写输出的地方
func (l *Logger) Write(p []byte) (n int, err error) {
	//加锁
	l.mu.Lock()
	//关锁
	defer l.mu.Unlock()
	//计算是否超过最大文件存储
	writeLen := int64(len(p))
	if writeLen > l.max() {
		return 0, fmt.Errorf(
			"write length %d exceeds maximum file size %d", writeLen, l.max(),
		)
	}
	//文件是否打开//如果用户没有传入则自己内部创建一个
	if l.file == nil {
		if err = l.openExistingOrNew(len(p)); err != nil {
			return 0, err
		}
	}

	if l.size+writeLen > l.max() {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = l.file.Write(p)
	l.size += int64(n)

	return n, err
}

// max 计算是否超过一个文件的最大存储
func (l *Logger) max() int64 {
	if l.MaxSize == 0 {
		return int64(defaultMaxSize * megabyte)
	}
	return int64(l.MaxSize) * int64(megabyte)
}

// openExistingOrNew 打开一个文件
func (l *Logger) openExistingOrNew(writeLen int) error {
	l.mill()

	//文件名字
	filename := l.filename()
	//返回文件具体信息
	info, err := os.Stat(filename)
	//检查是否因为目录是空引起的
	if os.IsNotExist(err) {
		return l.openNew()
	}
	if err != nil {
		return fmt.Errorf("error getting log file info: %s", err)
	}

	if info.Size()+int64(writeLen) >= l.max() {
		return l.rotate()
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// if we fail to open the old log file for some reason, just ignore
		// it and open a new log file.
		return l.openNew()
	}
	l.file = file
	l.size = info.Size()
	return nil
}

// genFilename generates the name of the logfile from the current time.
func (l *Logger) filename() string {
	if l.Filename != "" {
		return l.Filename
	}
	name := filepath.Base(os.Args[0]) + "-log_rotation.log"
	//这里，自己感觉linux可能会有问题//具体后边再说//如果有问题则改成下边
	// if l.Filename != "" {
	// 	return l.Filename
	// }
	// name := filepath.Base(os.Args[0]) + "-lumberjack.log"
	// return filepath.Join(os.TempDir(), name)
	str, _ := os.Getwd()
	return name + str
}

// openNew主要是用来 将已经不用的日志文件存放然后创建一个新的日志文件
func (l *Logger) openNew() error {
	//创建目录//如果没有这个目录
	err := os.MkdirAll(l.dir(), 0744)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	name := l.filename()
	mode := os.FileMode(0644)
	info, err := os.Stat(name)
	if err == nil {
		//获取权限
		mode = info.Mode()
		newname := backupName(name)
		//重命名
		if err := os.Rename(name, newname); err != nil {
			return fmt.Errorf("can't rename log file: %s", err)
		}

		if err := chown(name, info); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}
	l.file = f
	l.size = 0
	return nil
}

// dir() 获取文件中，目录部分
func (l *Logger) dir() string {
	return filepath.Dir(l.filename())
}

// backupName  组合形成新的名字
func backupName(name string) string {
	//获取name目录部分
	dir := filepath.Dir(name)
	//拿出基本文件名
	filename := filepath.Base(name)
	//获取扩展名
	ext := filepath.Ext(filename)
	//获取基本文件名的前缀部分，即不包括扩展名的部分。
	prefix := filename[:len(filename)-len(ext)]
	t := time.Now()
	//将时间转变成固定形式
	timestamp := t.Format(backupTimeFormat)
	//组合形成新的名字
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, timestamp, ext))
}

func (l *Logger) rotate() error {
	if err := l.close(); err != nil {
		return err
	}
	if err := l.openNew(); err != nil {
		return err
	}
	l.mill()
	return nil
}

// 从写关闭文件函数
func (l *Logger) close() error {
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

// mill有节制的去实现文件的压缩删除
func (l *Logger) mill() {
	l.startMill.Do(func() {
		//创建一个bool通道
		l.millCh = make(chan bool, 1)
		go l.millRun()
	})
	select {
	case l.millCh <- true:
	default:
	}
}

func (l *Logger) millRun() {
	for range l.millCh {
		_ = l.millRunOnce()
	}
}

// millRunOnce 执行文件的压缩和删除
func (l *Logger) millRunOnce() error {
	//三个都满足则说明没有要压缩和删除的文件
	if l.MaxNumbers == 0 && l.MaxSave == 0 && !l.Compress {
		return nil
	}
	//获取文件名按照时间
	files, err := l.oldLogFiles()
	if err != nil {
		return err
	}

	var compress, remove []logInfo
	//查看是否达到文件最大数量这个if是达到
	if l.MaxNumbers > 0 && l.MaxNumbers < len(files) {
		preserved := make(map[string]bool)
		var remaining []logInfo
		for _, f := range files {
			//获取文件名
			fn := f.Name()
			//查看是否为压缩文件.gz
			_ = strings.TrimSuffix(fn, compressSuffix)
			preserved[fn] = true

			if len(preserved) > l.MaxNumbers {
				remove = append(remove, f)
			} else {
				remaining = append(remaining, f)
			}
		}
		files = remaining
	}
	//查看是否为到期时间
	if l.MaxSave > 0 {
		diff := time.Duration(int64(24*time.Hour) * int64(l.MaxSave))
		cutoff := time.Now().Add(-1 * diff)

		var remaining []logInfo
		for _, f := range files {
			if f.timestamp.Before(cutoff) {
				remove = append(remove, f)
			} else {
				remaining = append(remaining, f)
			}
		}
		files = remaining
	}
	//查看是否压缩
	if l.Compress {
		for _, f := range files {
			//查看是否
			if !strings.HasSuffix(f.Name(), compressSuffix) {
				compress = append(compress, f)
			}
		}
	}
	//把所有要删除的删除
	for _, f := range remove {
		errRemove := os.Remove(filepath.Join(l.dir(), f.Name()))
		if err == nil && errRemove != nil {
			err = errRemove
		}
	}
	//未进行压缩的进行压缩
	for _, f := range compress {
		fn := filepath.Join(l.dir(), f.Name())
		errCompress := compressLogFile(fn, fn+compressSuffix)
		if err == nil && errCompress != nil {
			err = errCompress
		}
	}

	return err
}

// oldLogFiles 返回与当前日志文件同一目录中存储的备份日志文件列表
func (l *Logger) oldLogFiles() ([]logInfo, error) {
	//获取目录中所有日志文件
	files, err := os.ReadDir(l.dir())
	if err != nil {
		return nil, fmt.Errorf("can't read log file directory: %s", err)
	}
	//实体化一个对象
	logFiles := []logInfo{}
	prefix, ext := l.prefixAndExt()

	//遍历所有目录
	for _, f := range files {
		//IsDir()检验是否是一个目录
		if f.IsDir() {
			continue
		}
		if t, err := l.timeFromName(f.Name(), prefix, ext); err == nil {
			logFiles = append(logFiles, logInfo{t, f})
			continue
		}
		if t, err := l.timeFromName(f.Name(), prefix, ext+compressSuffix); err == nil {
			logFiles = append(logFiles, logInfo{t, f})
			continue
		}
		// error parsing means that the suffix at the end was not generated
		// by lumberjack, and therefore it's not a backup file.
	}
	//对切片进行排序
	sort.Sort(byFormatTime(logFiles))

	return logFiles, nil
}

// logInfo  文件的文件名和创建时间
type logInfo struct {
	timestamp time.Time
	os.DirEntry
}

// prefixAndExt 获取扩展名字和不带扩展的文件名
func (l *Logger) prefixAndExt() (prefix, ext string) {
	//获取文件名
	filename := filepath.Base(l.filename())
	//获取扩展名
	ext = filepath.Ext(filename)
	//获得不带扩展的文件名
	prefix = filename[:len(filename)-len(ext)] + "-"
	return prefix, ext
}

// timeFromName extracts the formatted time from the filename by stripping off
// the filename's prefix and extension. This prevents someone's filename from
// confusing time.parse.
func (l *Logger) timeFromName(filename, prefix, ext string) (time.Time, error) {
	//前缀检查
	if !strings.HasPrefix(filename, prefix) {
		return time.Time{}, errors.New("mismatched prefix")
	}
	//扩展名检查
	if !strings.HasSuffix(filename, ext) {
		return time.Time{}, errors.New("mismatched extension")
	}
	//提取时间字符串
	ts := filename[len(prefix) : len(filename)-len(ext)]
	//解析时间字符串
	return time.Parse(backupTimeFormat, ts)
}

// byFormatTime sorts by newest time formatted in the name.
type byFormatTime []logInfo

func (b byFormatTime) Less(i, j int) bool {
	return b[i].timestamp.After(b[j].timestamp)
}

func (b byFormatTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byFormatTime) Len() int {
	return len(b)
}

// compressLogFile 压缩文件
func compressLogFile(src, dst string) (err error) {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer f.Close()

	fi, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat log file: %v", err)
	}

	if err := chown(dst, fi); err != nil {
		return fmt.Errorf("failed to chown compressed log file: %v", err)
	}

	//创建一个压缩文件进行写操作
	gzf, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fi.Mode())
	if err != nil {
		return fmt.Errorf("failed to open compressed log file: %v", err)
	}

	defer gzf.Close()

	//创建一个新的写入器
	gz := gzip.NewWriter(gzf)

	//如果出现错误则删除创建的压缩文件
	defer func() {
		if err != nil {
			os.Remove(dst)
			err = fmt.Errorf("failed to compress log file: %v", err)
		}
	}()

	//下边进行拷贝和删除文件
	if _, err := io.Copy(gz, f); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	if err := gzf.Close(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Remove(src); err != nil {
		return err
	}

	return nil
}

// Close 线程安全的关闭文件
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.close()
}

func (l *Logger) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rotate()
}

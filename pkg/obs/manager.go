package obs

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/andreykaipov/goobs"
	"github.com/pkg/errors"
)

var (
	// Obs 路径
	ObsPath       = `C:\Program Files\obs-studio\bin\64bit\obs64.exe`
	ObsArg        = []string{"--profile", "zixun", "--scene", "zixun", "--minimize-to-tray", "--startvirtualcam"}
	ObsConfigPath = fmt.Sprintf(`%s\AppData\Roaming\obs-studio`, os.Getenv("USERPROFILE"))
)

type OBSManager struct {
	Client  *goobs.Client
	Cmd     *exec.Cmd
	Source  *VideoSource
	Profile *ProfileConfig
}

type VideoSource struct {
	Size *VideoSize
	Url  string
}

// OBSConfig OBS 配置
type OBSConfig struct {
	Size      *VideoSize
	StreamUrl string
	Mic       *VirtualMicInfo
}

func NewObsManager(source *VideoSource, port uint, password string) *OBSManager {
	salt := GenerateSalt()
	return &OBSManager{
		Source: source,
		Profile: &ProfileConfig{
			Size:     source.Size,
			Salt:     salt,
			Secret:   GenerateSecret(password, salt),
			Port:     port,
			Password: password,
			Mic: NewVirtualMicInfo(
				"CABLE Input (VB-Audio Virtual Cable)",
				"{0.0.0.00000000}.{ce4e42b4-c623-41d1-938f-83652535c9d0}"),
		},
	}
}

func (o *OBSManager) StopOBS() error {
	if o.Cmd != nil {
		log.Println("kill OBS")
		return o.Cmd.Process.Kill()
	}
	return nil
}

// cleanAllConfig 删除所有配置
func (o *OBSManager) cleanAllConfig() error {
	if err := os.RemoveAll(ObsConfigPath); err != nil {
		errors.Wrap(err, "delete config error")
	}

	p := strings.Join(
		[]string{
			ObsConfigPath, "global.ini",
		}, string(os.PathSeparator))

	if err := os.RemoveAll(p); err != nil {
		errors.Wrap(err, "delete config error")
	}
	return nil
}

// writeGlobalConfig 写入Global
//  @receiver o
//  @return error
func (o *OBSManager) writeGlobalConfig() error {
	c := GetGlobalConfigContent()

	p := strings.Join(
		[]string{
			ObsConfigPath, "global.ini",
		}, string(os.PathSeparator))
	if err := os.MkdirAll(filepath.Dir(p), 0644); err != nil {
		return errors.Wrap(err, "mkdir failed")
	}
	return os.WriteFile(p, []byte(c), 0644)

}

// writeProfileConfig 写入Profile
func (o *OBSManager) writeProfileConfig() error {
	c := GetProfileContent(o.Profile)

	p := strings.Join(
		[]string{
			ObsConfigPath, "basic", "profiles", "zixun", "basic.ini",
		}, string(os.PathSeparator))
	if err := os.MkdirAll(filepath.Dir(p), 0644); err != nil {
		return errors.Wrap(err, "mkdir failed")
	}
	return os.WriteFile(p, []byte(c), 0644)
}

// writeSceneConfig 写入Profile
func (o *OBSManager) writeSceneConfig() error {
	c := GetSceneContent(o.Source)

	p := strings.Join(
		[]string{
			ObsConfigPath, "basic", "scenes", "zixun.json",
		}, string(os.PathSeparator))
	if err := os.MkdirAll(filepath.Dir(p), 0644); err != nil {
		return errors.Wrap(err, "mkdir failed")
	}
	return os.WriteFile(p, []byte(c), 0644)

}

func (o *OBSManager) resetConfig() error {

	if err := o.cleanAllConfig(); err != nil {
		return err
	}

	if err := o.writeGlobalConfig(); err != nil {
		return err
	}

	if err := o.writeProfileConfig(); err != nil {
		return err
	}

	if err := o.writeSceneConfig(); err != nil {
		return err
	}

	return nil
}

func (o *OBSManager) StartOBS() error {

	if err := o.resetConfig(); err != nil {
		return err
	}

	cmd := exec.Command(ObsPath, ObsArg...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = filepath.Dir(ObsPath)

	log.Printf("run %s in path %s whit arg %s ", cmd.Path, cmd.Dir, cmd.Args)

	// cmd.Run()
	if err := cmd.Start(); err != nil {
		return err
	}
	log.Println(cmd.SysProcAttr)

	o.Cmd = cmd
	// time.Sleep(10 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := o.waitForClient(ctx)
	if err != nil {
		return err
	}
	o.Client = client
	return nil
}

func (o *OBSManager) waitForClient(ctx context.Context) (*goobs.Client, error) {
	for {
		client, err := goobs.New(
			fmt.Sprintf("127.0.0.1:%d", o.Profile.Port),
			goobs.WithPassword(o.Profile.Password), // optional
			goobs.WithDebug(true),                  // optional
		)
		if err != nil {
			select {
			case <-ctx.Done():
				return client, err
			default:
				log.Println("wait for client")
				time.Sleep(1 * time.Second)
				continue
			}
		}
		return client, err
	}
}

var allowedChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
var now = time.Now

type salt = string

func GenerateSalt() salt {
	r := rand.New(rand.NewSource(now().UnixNano()))
	sb := strings.Builder{}

	for i := 0; i < 32; i++ {
		sb.WriteByte(allowedChars[r.Intn(len(allowedChars))])
	}

	return base64.StdEncoding.EncodeToString([]byte(sb.String()))
}

// GenerateSecret 新建Secret
//  @param password
//  @param salt
//  @return string
func GenerateSecret(password string, salt salt) string {
	h := sha256.New()
	h.Write([]byte(password))
	h.Write([]byte(salt))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

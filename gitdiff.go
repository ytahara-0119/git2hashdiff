package main
import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)
func main() {
	// 引数の数をチェック
	if len(os.Args) != 1 {
		fmt.Println("エラー：引数は受け付けません")
		fmt.Println("使用法：git log --oneline | fzf --multi | gitdiff")
		os.Exit(1)
	}
	// 標準入力からハッシュ値を取得
	scanner := bufio.NewScanner(os.Stdin)
	var hashes []string
	for scanner.Scan() {
		hash := strings.Fields(scanner.Text())[0]
		hashes = append(hashes, hash)
	}
	// ハッシュ値が2つでない場合はエラー
	if len(hashes) != 2 {
		fmt.Println("エラー：2つのハッシュ値を指定してください")
		os.Exit(1)
	}
	// Git diffコマンドの実行
	cmd := exec.Command("git", "diff", hashes[0], hashes[1])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// fzfが終了するまで待機
	if err := cmd.Start(); err != nil {
		fmt.Println("Git diffコマンドの実行中にエラーが発生しました:", err)
		os.Exit(1)
	}
	// fzfプロセスの終了を待機
	err := cmd.Wait()
	if err != nil {
		// ExitErrorの場合、エラーコードが1のためエラーを表示しない
		if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.Sys().(syscall.WaitStatus).ExitStatus() != 1 {
			fmt.Println("プロセスを終了しました。")
		}
		os.Exit(1)
	}
}
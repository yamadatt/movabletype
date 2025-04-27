// movabletypeパッケージは「Movable Type インポート／エクスポート形式」のパース機能を提供します。
// Movable Typeのテキストベースのエクスポート形式をGoの構造体に変換できます。
package movabletype

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// デフォルト値
const (
	// AllowCommentsフィールドのデフォルト値（-1）
	DefaultAllowComments = -1

	// AllowPingsフィールドのデフォルト値（-1）
	DefaultAllowPings = -1
)

// EntryはMovable Typeインポート形式の1件の記事を表します。
// 各記事はMT仕様に基づくメタデータや本文などのフィールドを持ちます。
type Entry struct {
	Author   string // 著者名
	Title    string // 記事タイトル
	Basename string // URL用のベース名
	Status   string // 公開状態: "Draft", "Publish", "Future"

	// AllowCommentsはコメント許可設定（0または1）。未設定時はDefaultAllowComments。
	AllowComments int

	// AllowPingsはトラックバック/ピンバック許可設定（0または1）。未設定時はDefaultAllowPings。
	AllowPings int

	Converts        string    // 改行変換設定
	Date            time.Time // 公開日時
	PrimaryCategory string    // メインカテゴリ名
	Category        []string  // カテゴリ一覧
	Image           string    // アイキャッチ画像パス
	Body            string    // 本文
	ExtendedBody    string    // 追記本文
	Excerpt         string    // 抜粋・概要
	Keywords        string    // SEOキーワード
	Comment         string    // 記事へのコメント
}

// 新しいEntryをデフォルト値で生成します。
func NewEntry() *Entry {
	return &Entry{
		AllowComments: DefaultAllowComments,
		AllowPings:    DefaultAllowPings,
	}
}

// ParseはMovable Type形式のデータをio.Readerから読み込み、Entry構造体のスライスとして返します。
// 入力が不正な場合や必須フィールドに不正値がある場合はエラーを返します。
//
// 使用例:
//
//	entries, err := movabletype.Parse(os.Stdin)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, entry := range entries {
//	    fmt.Println(entry.Title)
//	}
func Parse(r io.Reader) ([]*Entry, error) {
	mts := []*Entry{}
	scanner := bufio.NewScanner(r)
	var err error
	m := NewEntry()

	for scanner.Scan() {
		line := scanner.Text()

		// 区切り線の処理
		if line == "-----" {
			continue
		}
		if line == "--------" {
			mts = append(mts, m)
			m = NewEntry()
			continue
		}

		// 複数行フィールドの処理
		if strings.HasSuffix(line, ":") {
			field := line[:len(line)-1] // ":"を除去
			content := ""
			for scanner.Scan() {
				l := scanner.Text()
				if l == "-----" {
					break
				}
				content += l + "\n"
			}
			switch field {
			case "BODY":
				m.Body = content
			case "EXTENDED BODY":
				m.ExtendedBody = content
			case "EXCERPT":
				m.Excerpt = content
			case "KEYWORDS":
				m.Keywords = content
			case "COMMENT":
				m.Comment = content
			}
			continue
		}

		// 1行フィールドの処理
		ss := strings.SplitN(line, ": ", 2)
		if len(ss) != 2 {
			continue
		}
		key, value := ss[0], ss[1]
		switch key {
		case "AUTHOR":
			m.Author = value
		case "TITLE":
			m.Title = value
		case "BASENAME":
			m.Basename = value
		case "STATUS":
			if value == "Draft" || value == "Publish" || value == "Future" {
				m.Status = value
			} else {
				return nil, fmt.Errorf("STATUS列はDraft, Publish, Futureのみ許可されています。取得値: %s", value)
			}
		case "ALLOW COMMENTS":
			m.AllowComments, err = strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("ALLOW COMMENTS列は0または1のみ許可: %w", err)
			}
			if m.AllowComments != 0 && m.AllowComments != 1 {
				return nil, fmt.Errorf("ALLOW COMMENTS列は0または1のみ許可。取得値: %d", m.AllowComments)
			}
		case "ALLOW PINGS":
			m.AllowPings, err = strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("ALLOW PINGS列は0または1のみ許可: %w", err)
			}
			if m.AllowPings != 0 && m.AllowPings != 1 {
				return nil, fmt.Errorf("ALLOW PINGS列は0または1のみ許可。取得値: %d", m.AllowPings)
			}
		case "CONVERT S":
			m.Converts = value
		case "DATE":
			if strings.HasSuffix(value, "AM") || strings.HasSuffix(value, "PM") {
				m.Date, err = time.Parse("01/02/2006 03:04:05 PM", value)
			} else {
				m.Date, err = time.Parse("01/02/2006 15:04:05", value)
			}
			if err != nil {
				return nil, fmt.Errorf("DATE列のパースエラー: %w", err)
			}
		case "PRIMARY CATEGORY":
			m.PrimaryCategory = value
		case "CATEGORY":
			m.Category = append(m.Category, value)
		case "IMAGE":
			m.Image = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return mts, nil
}

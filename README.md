# movabletype

movabletypeパッケージは「Movable Type インポート／エクスポート形式」のテキストデータをパースし、Goの構造体として扱えるようにするライブラリです。

## 概要

- Movable Typeのエクスポートファイル（ブログ記事やコメントなどが記載されたテキスト形式）を読み込み、各記事をGoの構造体（Entry）として取得できます。
- 各記事（Entry）は、著者名・タイトル・本文・カテゴリ・公開状態・コメント許可設定など、Movable Typeの仕様に沿ったフィールドを持っています。
- 標準入力やファイルからデータを読み込んで、記事ごとに分割し、各種フィールドを自動的にパースします。

## サンプルコード

```go
package main

import (
	"fmt"
	"os"

	"github.com/yamadatt/movabletype"
)

func main() {
	entries, err := movabletype.Parse(os.Stdin)
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		fmt.Printf("タイトル: %s\n", e.Title)
		fmt.Printf("本文: %s\n", e.Body)
	}
}
```

## Entry構造体の主なフィールド

| フィールド名         | 型         | 説明                         |
|----------------------|------------|------------------------------|
| Author               | string     | 著者名                       |
| Title                | string     | 記事タイトル                 |
| Basename             | string     | URL用のベース名              |
| Status               | string     | 公開状態: "Draft", "Publish", "Future" |
| AllowComments        | int        | コメント許可（0または1）      |
| AllowPings           | int        | ピンバック許可（0または1）    |
| Converts             | string     | 改行変換設定                 |
| Date                 | time.Time  | 公開日時                     |
| PrimaryCategory      | string     | メインカテゴリ名             |
| Category             | []string   | カテゴリ一覧                 |
| Image                | string     | アイキャッチ画像パス         |
| Body                 | string     | 本文                         |
| ExtendedBody         | string     | 追記本文                     |
| Excerpt              | string     | 抜粋・概要                   |
| Keywords             | string     | SEOキーワード                |
| Comment              | string     | 記事へのコメント             |

## 注意事項
- Movable Typeのエクスポート形式に厳密に従っていないデータの場合、パース時にエラーとなる場合があります。
- 日付のパースは「01/02/2006 15:04:05」または「01/02/2006 03:04:05 PM」形式に対応しています。

---

ご不明点や追加の使い方例が必要な場合は、IssueやPRでご連絡ください。


## 参考リンク

- [MovableType.org – Documentation: The Movable Type Import / Export Format](https://movabletype.org/documentation/appendices/import-export-format.html)
- (日本語) [記事のインポートフォーマット : Movable Type 6 ドキュメント](https://www.movabletype.jp/documentation/mt6/tools/import-export-format.html)
# chat-app_fg
Chat app with golang and flutter
Go を 「Goプログラミング実践入門」(インプレス)で勉強中

fix：バグ修正
add：新規（ファイル）機能追加
update：機能修正（バグではない）
remove：削除（ファイル）

feature ブランチで作業し main に統合する

# Arch
yay -S postgresql
sudo -iu postgres
su
passwd postgres //パスワード設定
su postgres
initdb -D /var/lib/postgres/data
systemctl start postgresql
createuser --interactive //ユーザー作成、自分と同じ名前にすると楽
createdb ユーザー名
メインユーザーで
createdb chat
psql -f setup.sql -d chat
//データベースの権限周りが全然わからん 適当

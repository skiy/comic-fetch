<?php

define('DEBUG', true);
define('BASEPATH', __DIR__ . DIRECTORY_SEPARATOR);

ini_set("display_errors", DEBUG);
error_reporting(E_ALL);

$bookId = $_GET['book_id'] ?? '';
$orderBy = $_GET['order'] ?? '0';

if ($bookId <= 0) {
    exit('漫画编号参数错误');
}

include 'db.php';

$bookSql = "SELECT * FROM tb_books WHERE id = :id";
$sth = $dbh->prepare($bookSql);
$sth->execute(array(':id' => $bookId));
$bookInfo = $sth->fetch();
?>
<!doctype html>
<html lang="en">
<head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->

    <title>第五城市漫画网</title>
</head>
<body>
<?php
if ($bookInfo === false) {
    exit('漫画不存在');
}

$orderBy = ($orderBy == 0) ? 'DESC' : 'ASC';
$chapterSql = "SELECT * FROM tb_chapters WHERE bid = :bid ORDER BY id {$orderBy}";
$sth = $dbh->prepare($chapterSql);
$sth->execute(array(':bid' => $bookId));
$chapterInfo = $sth->fetchAll();

if (empty($chapterInfo)) {
    exit("<<{$bookInfo['name']}>> 漫画不存在章节");
}

foreach ($chapterInfo as $key => $value) {
    $url = "chapter.php?book_id={$value['bid']}&chapter_id={$value['id']}";
    echo "<a href='{$url}' target='_blank'>{$value['title']}</a><br />";
}
?>
</body>
</html>
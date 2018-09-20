<?php

define('DEBUG', true);

ini_set("display_errors", DEBUG);
error_reporting(E_ALL);

$bookId = $_GET['book_id'] ?? '';
$orderBy = $_GET['order'] ?? '0';

if ($bookId <= 0) {
    exit('漫画编号参数错误');
}

$dbh = new PDO('mysql:host=127.0.0.1;dbname=comic', 'root', '123456');

$bookSql = "SELECT * FROM tb_books WHERE id = :id";
$sth = $dbh->prepare($bookSql);
$sth->execute(array(':id' => $bookId));
$bookInfo = $sth->fetch();

if ($bookInfo === false) {
    exit('漫画不存在');
}

$orderBy = ($orderBy == 0) ? 'DESC' : 'ASC';
$chapterSql = "SELECT * FROM tb_chapters WHERE bid = :bid ORDER BY chapter_id {$orderBy}";
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
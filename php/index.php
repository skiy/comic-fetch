<?php

/**
 *
 * File:  index.php
 * Author: Skiychan <dev@skiy.net>
 * Created: 2018/09/22
 */

define('DEBUG', true);
define('BASEPATH', __DIR__ . DIRECTORY_SEPARATOR);

ini_set("display_errors", DEBUG);
error_reporting(E_ALL);

include BASEPATH . 'db.php';

$bookSql = "SELECT * FROM tb_books ORDER BY id DESC ";
$sth = $dbh->prepare($bookSql);
$sth->execute();
$bookInfo = $sth->fetchAll();

if (empty($bookInfo)) {
    exit("未抓取到漫画");
}

foreach ($bookInfo as $key => $value) {
    $url = "comic.php?book_id={$value['id']}";
    echo "<a href='{$url}' target='_blank'>{$value['name']}</a><br />";
}
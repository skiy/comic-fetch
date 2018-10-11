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
if (empty($bookInfo)) {
    exit("未抓取到漫画");
}

foreach ($bookInfo as $key => $value) {
    $url = "comic.php?book_id={$value['id']}";
    if (! empty($value['image_url'])) {
        echo "<img src='images/{$value['image_url']}' /><br />";
    }
    echo "<a href='{$url}' target='_blank'>{$value['name']}</a><br />";
}
?>
</body>
</html>
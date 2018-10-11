<?php

define('DEBUG', true);
define('BASEPATH', __DIR__ . DIRECTORY_SEPARATOR);

ini_set("display_errors", DEBUG);
error_reporting(E_ALL);

$bookId = $_GET['book_id'] ?? '';
$chapterId = $_GET['chapter_id'] ?? '';

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
if ($bookId <= 0 || $chapterId <= 0) {
    exit('漫画章节参数错误');
}

include BASEPATH . 'db.php';

$imageSql = "SELECT * FROM tb_images WHERE bid = :bid AND cid = :cid ORDER BY id ASC";
$sth = $dbh->prepare($imageSql);
$sth->execute(array(':bid' => $bookId, ':cid' => $chapterId));
$imageInfo = $sth->fetchAll();

if (empty($imageInfo)) {
    exit("本话漫画不存在");
}

foreach ($imageInfo as $key => $value) {
    $url = $value['image_url'];
    if (empty($value['image_url'])) {
        $url = 'img.php?url=' . $value['origin_url'];
    } else {
        $url = "images/{$value['image_url']}";
    }
    echo "<img src='{$url}' /><br />";
}
?>
</body>
</html>
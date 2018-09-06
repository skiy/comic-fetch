<?php

/**
 * GET 请求
 * @param string $url 链接
 * @param array $headers 用户头部信息
 * @return string content
 */
function http_get($url, $headers = array()) {
    $oCurl = curl_init();
    if (stripos($url, "https://") !== FALSE) {
        curl_setopt($oCurl, CURLOPT_SSL_VERIFYPEER, FALSE);
        curl_setopt($oCurl, CURLOPT_SSL_VERIFYHOST, FALSE);
        curl_setopt($oCurl, CURLOPT_SSLVERSION, 1); //CURL_SSLVERSION_TLSv1
    }
    curl_setopt($oCurl, CURLOPT_URL, $url);
    curl_setopt($oCurl, CURLOPT_RETURNTRANSFER, 1);
    curl_setopt($oCurl, CURLOPT_HTTPHEADER, $headers);
    $sContent = curl_exec($oCurl);
    $aStatus = curl_getinfo($oCurl);
    curl_close($oCurl);
    if (intval($aStatus["http_code"]) == 200) {
        return $sContent;
    } else {
        return false;
    }
}

$url = $_GET['url'] ?? exit('图片链接错误');

$headers = array(
    'DNT: 1',
    'If-Modified-Since: Thu, 06 Sep 2018 03:54:19 GMT',
    'If-None-Match: "BDE9E8B0317BF99A37BE8FE52763AF1E"',
    'Referer: https://www.mh160.com',
);

header('content-type: image/jpeg');
echo http_get($url, $headers);
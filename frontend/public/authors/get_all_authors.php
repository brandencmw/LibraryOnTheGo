<?php
function getAuthors(string $order='', $name='') : array {

    $get_all_endpoint = 'https://server/authors?page_size=20&include_images=true';
    if ($order != '') {
        $get_all_endpoint .= '&order=' . $order;
    }

    if ($name != '') {
        $get_all_endpoint .= '&name=' . $name;
    }

    $ch = curl_init($get_all_endpoint);
    curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'GET');
    curl_setopt($ch, CURLOPT_SSLVERSION, CURL_SSLVERSION_TLSv1_3);
    curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, true);
    curl_setopt($ch, CURLOPT_CAINFO, '/etc/ssl/certs/root-ca.pem');
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);

    $json_response = curl_exec($ch);
    $authors = [];
    if (curl_errno($ch)) {
        echo curl_error($ch);
        return $authors;
    }
    curl_close($ch);
    $response = json_decode($json_response, true);

    return $response['authors'];
}
?>
<?php
function getAllBooks() {
    $get_books_endpoint = 'https://server/books';

    $ch = curl_init($get_books_endpoint);
    curl_setopt($ch, CURLOPT_CUSTOMREQUEST,'GET');
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, true);
    curl_setopt($ch, CURLOPT_CAINFO, '/etc/ssl/certs/root-ca.pem');

    $json_response = curl_exec($ch);
    if (curl_errno($ch)) {
        echo ''. curl_error($ch);
    }
    curl_close($ch);

    $response = json_decode($json_response, true);
    return isset($response['books']) ? $response['books'] : [];
}
?>
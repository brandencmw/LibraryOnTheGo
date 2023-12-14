<?php
function searchBooks($title, $author) {
    $search_endpoint = 'https://google_books_service/book?title=' . $title . '&author=' . $author;

    $ch = curl_init($search_endpoint);
    curl_setopt($ch, CURLOPT_CUSTOMREQUEST,'GET');
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, true);
    curl_setopt($ch, CURLOPT_CAINFO, '/etc/ssl/certs/root-ca.pem');

    $json_response = curl_exec($ch);
    curl_close($ch);
    if (curl_errno($ch)) {
        echo json_encode(''. curl_error($ch));
    }
    curl_close($ch);
    return trim($json_response);
}


if (isset($_GET['title']) && isset($_GET['author'])) {
    $results = searchBooks($_GET['title'], $_GET['author']);
    echo $results;
}


?>
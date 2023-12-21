<?php
function deleteBook($authToken, $id) {

    $delete_endpoint = 'https://server/books/auth/delete?id=' . $id;
    $ch = curl_init($delete_endpoint);

    curl_setopt($ch, CURLOPT_CUSTOMREQUEST,'DELETE');
    curl_setopt($ch, CURLOPT_CAINFO, '/etc/ssl/certs/root-ca.pem');
    curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, true);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER,1);
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        'Authorization: Bearer ' . $authToken,
    ]);

    $res = curl_exec($ch);
    if (curl_errno($ch)) {
        echo "". curl_error($ch);
        exit();
    }
    echo $res;
}

if ($_SERVER['REQUEST_METHOD'] == 'DELETE') {
    if (isset($_COOKIE['Authorization'])) {
        deleteBook($_COOKIE['Authorization'], $_GET['id']);
    } else {
        header('Location: ../../../login/');
    }
}
?>
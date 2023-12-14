<?php declare(strict_types=1);

function authenticated(string $token) : bool {
    $authEndpoint = 'https://server/auth/';
    $ch = curl_init($authEndpoint);
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        "Authorization: Bearer " . $token,
    ]);
    curl_setopt($ch, CURLOPT_POST, true);
    curl_setopt($ch, CURLOPT_SSLVERSION, CURL_SSLVERSION_TLSv1_3);
    curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, true);
    curl_setopt($ch, CURLOPT_CAINFO, '/etc/ssl/certs/root-ca.pem');
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_VERBOSE, true);
    $response = curl_exec($ch);#
    if (curl_errno($ch)) {
        echo curl_error($ch);
        exit;
    }
    curl_close($ch);
    if ($response) {
        $json_response = json_decode($response, true);
        return isset($json_response['authenticated']) && $json_response['authenticated'];
    }
    return false;
}

?>
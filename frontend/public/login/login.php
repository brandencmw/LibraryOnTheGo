<?php
    if ($_SERVER['REQUEST_METHOD'] === 'POST') {
        $username = $_POST['username'];
        $password = $_POST['password'];

        $loginEndpoint = 'https://server/auth/login';

        $req = [
            'username' => $username,
            'password' => $password,
        ];
        $json_data = json_encode($req);
    
        $ch = curl_init($loginEndpoint);
        curl_setopt($ch, CURLOPT_POST, true);
        curl_setopt($ch, CURLOPT_SSLVERSION, CURL_SSLVERSION_TLSv1_3);
        curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, true);
        curl_setopt($ch, CURLOPT_CAINFO, '/etc/ssl/certs/root-ca.pem');
        curl_setopt($ch, CURLOPT_POSTFIELDS, $json_data);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_HEADER, true);
        curl_setopt($ch, CURLOPT_HTTPHEADER, [
            "Content-Type: application/json",
            "Content-Length: " . strlen($json_data),
        ]);

        $response = curl_exec($ch);
        if (curl_errno($ch)) {
            echo 'Curl error: ' . curl_error($ch);
            exit();
        }

        $header_size = curl_getinfo($ch, CURLINFO_HEADER_SIZE);
        $headers = substr($response,0, $header_size);
        if (preg_match('/Authorization=\s*(.*?)(?:\r\n|$|;)/', $headers, $matches)) {
            $authToken = $matches[1];
            $cookiePath = preg_match('/Path=\s*(.*?)(?:\r\n|$|;)/', $headers, $matches) ? $matches[1] :'/';
            $cookieDomain = preg_match('/Domain=\s*(.*?)(?:\r\n|$|;)/', $headers, $matches) ? $matches[1] : 'localhost';
            $cookieExpiry = preg_match('/Max-Age=\s*(.*?)(?:\r\n|$|;)/', $headers, $matches) ? (int)$matches[1] : 3600;
            
            // Set the cookie in the client's browser
            setcookie('Authorization', $authToken, [
                'expires' => time() + $cookieExpiry,
                'path' => $cookiePath,
                'domain' => $cookieDomain,
                'secure' => preg_match('/Secure\s*(.*?)(?:\r\n|$|;)/', $headers, $matches),
                'httponly' => preg_match('/HttpOnly\s*(.*?)(?:\r\n|$|;)/', $headers, $matches),
                'samesite' => 'none'
            ]);
            header('Location: ../admin');
        }
        curl_close($ch);
    } else {
        echo 'Invalid request.';
    }
?>
<?
class Author {
    public string $id;
    public string $firstName;
    public string $lastName;
    public string $bio;
    public string $imageReference;

    public function __construct(string $id, string $firstName, string $lastName, string $bio, string $imageReference) {
        $this->id = $id;
        $this->firstName = $firstName;
        $this->lastName = $lastName;
        $this->bio = $bio;
        $this->imageReference = $imageReference;
    }
}
function getAuthors() : array {

    $get_all_endpoint = 'https://server/authors';
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
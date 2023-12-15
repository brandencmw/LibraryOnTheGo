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


?>
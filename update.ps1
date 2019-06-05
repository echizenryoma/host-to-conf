if ($PSScriptRoot -ne "") {
    Set-Location $PSScriptRoot
}

Invoke-WebRequest -Uri "https://raw.githubusercontent.com/googlehosts/hosts/master/hosts-files/hosts" -OutFile hosts
& go build .
& .\host-to-conf

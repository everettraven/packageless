Write-Output "Creating any needed directories"

if (-not (Test-Path %USERPROFILE%\bin)) {
    New-Item "~/bin" -ItemType Directory
}

if (-not (Test-Path %USERPROFILE%\.packageless)) {
    New-Item %USERPROFILE%\.packageless -ItemType Directory
    New-Item %USERPROFILE%\.packageless\pims_config -ItemType Directory
    New-Item %USERPROFILE%\.packageless\pims -ItemType Directory 
}

Write-Output "Downloading the executable"

Invoke-WebRequest https://github.com/everettraven/packageless/releases/latest/download/packageless-windows.exe -OutFile %USERPROFILE%\bin\packageless.exe

Write-Output "Adding packageless to PATH"

setx PATH "%PATH%;%USERPROFILE%\bin\packageless"
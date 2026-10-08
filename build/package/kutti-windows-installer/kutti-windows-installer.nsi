!ifndef VERSION
    !define VERSION "0.0.0"
!endif
!ifndef VERSION_FOUR_PART
    !define VERSION_FOUR_PART "0.0.0.0"
!endif

Name "Kutti ${VERSION}"
OutFile "out/kutti-windows-installer.exe"
Unicode True
Icon "cmd/kutti/winres/icon.ico"
VIProductVersion "${VERSION_FOUR_PART}"
VIAddVersionKey "ProductName" "Kutti"
VIAddVersionKey "ProductVersion" "${VERSION}"
VIAddVersionKey "FileVersion" "${VERSION_FOUR_PART}"
VIAddVersionKey "FileDescription" "Kutti CLI Installer"
VIAddVersionKey "LegalCopyright" "Kutti Project"

RequestExecutionLevel user
InstallDir "$LOCALAPPDATA\Programs\kutti"

Section
    SetOutPath $INSTDIR

    # Create the uninstaller
    WriteUninstaller "$INSTDIR\uninstall.exe"

    # Files
    File /oname=kutti.exe out/kutti_windows_amd64.exe
    File /oname=kutti.ico cmd/kutti/winres/icon.ico
    File /oname=kutti-completion.ps1 out/completions/kutti-completion.ps1

    # Windows Terminal JSON Fragment
    SetOutPath "$LOCALAPPDATA\Microsoft\Windows Terminal\Fragments\Kutti\"
    File /oname=kutti.json build/package/kutti-windows-installer/kutti-wt-profile.json

    # Shortcuts
    CreateShortcut "$SMPROGRAMS\Uninstall Kutti ${VERSION}.lnk" "$INSTDIR\uninstall.exe"

    # Set Output path before creating shortcuts, as that will be
    # the working directory for the shortcuts. 
    SetOutPath "$DOCUMENTS"
    CreateShortcut "$SMPROGRAMS\Kutti Command Prompt.lnk" "%windir%\system32\cmd.exe" '/K "PATH=%PATH%;$INSTDIR"' "$INSTDIR\kutti.ico"
    CreateShortcut "$SMPROGRAMS\Kutti PowerShell.lnk" "%SystemRoot%\system32\WindowsPowerShell\v1.0\powershell.exe" `-NoExit -ExecutionPolicy Bypass -Command "$$env:Path += ';$INSTDIR'; if (Test-Path '$INSTDIR\kutti-completion.ps1') { . '$INSTDIR\kutti-completion.ps1' }"` "$INSTDIR\kutti.ico"

    # Register in Windows Installed Apps (Add/Remove Programs)
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\Kutti" "DisplayName" "Kutti ${VERSION}"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\Kutti" "DisplayVersion" "${VERSION}"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\Kutti" "Publisher" "Kutti Project"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\Kutti" "DisplayIcon" "$INSTDIR\kutti.ico"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\Kutti" "InstallLocation" "$INSTDIR"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\Kutti" "UninstallString" '"$INSTDIR\uninstall.exe"'
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\Kutti" "QuietUninstallString" '"$INSTDIR\uninstall.exe" /S'
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\Kutti" "URLInfoAbout" "https://github.com/kuttiproject/kutti"
    WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\Kutti" "NoModify" 1
    WriteRegDWORD HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\Kutti" "NoRepair" 1
SectionEnd

# Uninstaller section
Section "Uninstall"
    # Always delete uninstaller first
    Delete "$INSTDIR\uninstall.exe"

    # Delete files
    Delete "$INSTDIR\kutti.exe"
    Delete "$INSTDIR\kutti.ico"
    Delete "$INSTDIR\kutti-completion.ps1"

    # Remove Windows Terminal JSON Fragment
    RMDir /r "$LOCALAPPDATA\Microsoft\Windows Terminal\Fragments\Kutti\"

    # Delete shortcuts
    Delete "$SMPROGRAMS\Uninstall Kutti*.lnk"
    Delete "$SMPROGRAMS\Kutti Command Prompt.lnk"
    Delete "$SMPROGRAMS\Kutti PowerShell.lnk"

    # Delete Registry entry from Installed Apps
    DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\Kutti"

    # Delete the directory
    RMDir $INSTDIR
SectionEnd

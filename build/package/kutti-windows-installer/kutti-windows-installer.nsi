Name "Kutti Installer"
OutFile "out/kutti-windows-installer.exe"
Unicode True
Icon "cmd/kutti/winres/icon.ico"

RequestExecutionLevel user
InstallDir "$LOCALAPPDATA\Programs\kutti"

Section
    SetOutPath $INSTDIR

    # Create the uninstaller
    WriteUninstaller "$INSTDIR\uninstall.exe"

    # Files
    File /oname=kutti.exe out/kutti_windows_amd64.exe
    File /oname=kutti.ico cmd/kutti/winres/icon.ico

    # Windows Terminal JSON Fragment
    SetOutPath "$LOCALAPPDATA\Microsoft\Windows Terminal\Fragments\Kutti\"
    File /oname=kutti.json build/package/kutti-windows-installer/kutti-wt-profile.json

    # Shortcuts
    CreateShortcut "$SMPROGRAMS\Uninstall Kutti.lnk" "$INSTDIR\uninstall.exe"

    # Set Output path before creating shortcuts, as that will be
    # the working directory for the shortcuts. 
    SetOutPath "$DOCUMENTS"
    CreateShortcut "$SMPROGRAMS\Kutti Command Prompt.lnk" "%windir%\system32\cmd.exe" '/K "PATH=%PATH%;$INSTDIR"'
    CreateShortcut "$SMPROGRAMS\Kutti PowerShell.lnk" "%SystemRoot%\system32\WindowsPowerShell\v1.0\powershell.exe" `-NoExit -c "$$env:Path += ';$INSTDIR'"`
SectionEnd



# create a section to define what the uninstaller does.
# the section will always be named "Uninstall"
Section "Uninstall"
    # Always delete uninstaller first
    Delete "$INSTDIR\uninstall.exe"

    # Delete files
    Delete "$INSTDIR\kutti.exe"
    Delete "$INSTDIR\kutti.ico"

    # Remove Windows Terminal JSON Fragment
    RMDir /r "$LOCALAPPDATA\Microsoft\Windows Terminal\Fragments\Kutti\"

    # Delete shortcuts
    Delete "$SMPROGRAMS\Uninstall Kutti.lnk"
    Delete "$SMPROGRAMS\Kutti Command Prompt.lnk"
    Delete "$SMPROGRAMS\Kutti PowerShell.lnk"

    # Delete the directory
    RMDir $INSTDIR
SectionEnd

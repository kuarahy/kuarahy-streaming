!include "MUI2.nsh"

Name "Twitch Notifier"
OutFile "TwitchNotifierSetup.exe"
InstallDir "$PROGRAMFILES\TwitchNotifier"
RequestExecutionLevel admin

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "LICENSE"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_LANGUAGE "English"

Section "Install"
  SetOutPath $INSTDIR
  File "twitch-notifier.exe"
  File "config.default.yaml"
  
  # Create uninstaller
  WriteUninstaller "$INSTDIR\Uninstall.exe"
  
  # Create start menu shortcut
  CreateDirectory "$SMPROGRAMS\Twitch Notifier"
  CreateShortCut "$SMPROGRAMS\Twitch Notifier\Twitch Notifier.lnk" "$INSTDIR\twitch-notifier.exe"
  CreateShortCut "$SMPROGRAMS\Twitch Notifier\Uninstall.lnk" "$INSTDIR\Uninstall.exe"
  
  # Add to PATH
  EnVar::AddValue "PATH" "$INSTDIR"
  
  # Register as a service (optional)
  ExecWait '"$INSTDIR\twitch-notifier.exe" install'
SectionEnd

Section "Uninstall"
  # Stop service if running
  ExecWait '"$INSTDIR\twitch-notifier.exe" uninstall'
  
  # Remove files
  Delete "$INSTDIR\twitch-notifier.exe"
  Delete "$INSTDIR\config.yaml"
  Delete "$INSTDIR\Uninstall.exe"
  
  # Remove shortcuts
  Delete "$SMPROGRAMS\Twitch Notifier\Twitch Notifier.lnk"
  Delete "$SMPROGRAMS\Twitch Notifier\Uninstall.lnk"
  RMDir "$SMPROGRAMS\Twitch Notifier"
  
  # Remove from PATH
  EnVar::DeleteValue "PATH" "$INSTDIR"
  
  # Remove installation directory
  RMDir /r "$INSTDIR"
SectionEnd
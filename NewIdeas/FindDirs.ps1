# Powershell script to parse out log file locations for IIS
# Run with powershell -executionpolicy bypass <Location of this script> <Directory>
$DirName=$args[0]
# Copy the IIS XML File to <Directory> so we have our own copy in the collection
$SysRoot = $Env:SYSTEMROOT
Copy-Item "$SysRoot\System32\inetsrv\config\applicationHost.config" -Destination "$DirName"
# Start with some standard Locations
Add-Content -Path "$DirName\DirOutput.txt" -Value "%SystemRoot%\System32\LogFiles"
Add-Content -Path "$DirName\DirOutput.txt" -Value "%SystemDrive%\Windows.old\Windows\System32\LogFiles"
Add-Content -Path "$DirName\DirOutput.txt" -Value "%SystemDrive%\inetpub\logs\LogFiles"
# Use Powershell to parse the XPaths of the Log Locations
Select-Xml -Path $DirName\applicationHost.config -XPath '/configuration/system.applicationHost/log/centralBinaryLogFile' | ForEach-Object { $_.Node.directory } | Add-Content -Path "$DirName\DirOutput.txt"
Select-Xml -Path $DirName\applicationHost.config -XPath '/configuration/system.applicationHost/log/centralW3CLogFile' | ForEach-Object { $_.Node.directory } | Add-Content -Path "$DirName\DirOutput.txt"
Select-Xml -Path $DirName\applicationHost.config -XPath '/configuration/system.applicationHost/sites/site/logFile' | ForEach-Object { $_.Node.directory } | Add-Content -Path "$DirName\DirOutput.txt"
Select-Xml -Path $DirName\applicationHost.config -XPath '/configuration/system.applicationHost/sites/siteDefaults/logFile' | ForEach-Object { $_.Node.directory } | Add-Content -Path "$DirName\DirOutput.txt"
Select-Xml -Path $DirName\applicationHost.config -XPath '/configuration/system.applicationHost/sites/siteDefaults/traceFailedRequestsLogging' | ForEach-Object { $_.Node.directory } | Add-Content -Path "$DirName\DirOutput.txt"
# Now Dedup The File
gc "$DirName\DirOutput.txt" | sort | get-unique > "$DirName\DirOutputDedup.txt"

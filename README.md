# AChoirX
 AChoirX is a port of AChoir in Go(lang) for Cross Platform forensic collection, processing, and analysis.

The goal of AChoirX is to create a simple scripting language for Targeted Collection (Live Response) on Windows, Linux, OSX (MacOS), and Android.  AChoirX (like AChoir) will also support scripting of extraction, parsing, and analysis of mounted forensic images.

The AChoirX installation contains a single AChoir collection script, which can run on all three primary versions of AChoirX (Windows, Linux, MacOS). This shows the power of the idea: A single consistent method to gather artifacts and telemetry from all three platforms using a single script.  AChoirX does not to abstract these artifacts - but rather collects the actual platform specific raw artifacts in a consistent manner using the same utility on all three platforms.  Please Note: While AChoirX also runs on Android - it is not yet supported in the default script.

AChoir scripts can be modified in nearly endless ways to meet the collection needs of different types of investigations.  AChoirX can also run in interactive console mode, or across TCP using the built in Multi-Handler.  AChoirX also natively supports common upload methods such as SFTP or S3.  Since AChoirX runs on Windows, Linux, MacOS, and Android - All of these functions, including the multi-handler (client and server modes) and the upload methods, work regardless of platform.

# Change Log:
10/16/2020 - Alpha 23. Windows, Linux, OSX supported. The default script runs on all three platforms and gathers platform specific Artifacts and Telemetry. Has built in S3 Bucket uploading (using the AWS SDK).

10/17/2020 - Alpha 24. Hash Running Program for non-repudiation. Add &Myp (My Program) and &Myh (My Hash)

10/25/2020 - Alpha 25. Add /GXR: - Gets a Zip File, Extracts it, and runs the script.

10/30/2020 - Alpha 26. Copy Running Program to \Progs for non-repudiation

11/07/2020 - Beta 27. Upgraded Status to Beta.  Change CPY: Target File Atime and MTime to match Source. Change FileExists to accept File or Directory & Improve Error handling. Add Quoted Parsing to EXE: and SYS: processing.

11/10/2020 - Beta 28. Add &USR and &PWD - To enable UserID and Password on Command Line

11/13/2020 - Beta 29. Implement Encryption (/DEC:, ENC:, DEC:).

11/18/2020 - Beta 30. Minor Fixes.  Added PlasoX.ACQ Script for cross platform Plaso (Targeted TimeLine) processing

11/26/2020 - Beta 31. Add Admin/Root Checks - Move subroutine to Platform Specific Files.  Add &Adm Variable = Yes or No  (Running as Admin/Root)

11/27/2020 - Beta 32. Add &Mem Variable (Total System Memory)

11/29/2020 - Beta 33. Add support for unzipping an embedded []byte stream. Embed Platform specific Default Scripts (Win, Lin, OSX)

12/20/2020 - Beta 34. Add Embedder to ToolChain.  Include WinPmem (Memory Dumper) in Embedded Zip. Add TSK Fcat into (Raw NTFS Copy) into Embedded Zip

12/20/2020 - Beta 35. Fix &CNR Counter

12/27/2020 - Beta 36. Change Conditional Logic to only count a single occurance of &FOR and &LST comparisons. This prevents the need for multiple END: statements.
Multiple comparisons only get a single hit if ANY match is found. THIS IS IMPORTANT!! Wherever &FOR and &LST are used in CONDITIONAL LOGIC - A SINGLE HIT WILL BE 
TRUE.  To Test for INDIVIDUAL cases use a specific check and NOT a Check Against a list (&LST, &FOR). Expand &FOR and &LST Support to more Actions. Add 
HSH:<Filename> Will put the File hash in the &HSH Variable (Only supports a single File for now). Trim quotes for CKN: and CKY:

12/28/2020 - Beta 37 - Implement CopyPath= for Single File Copy

01/01/2021 - Beta 38 - Implement END:Reset to clear any Dangling ENDs.  Use Judiciously.

02/14/2021 - Beta 39 - Implement native SFTP Upload

02/21/2021 - Beta 40(RC1) - Refactor code for efficiencies, Fix Command Line Variables. Improve Comparisons for Missing Parameters (EQU:, NEQ:, N==:, N<<:, N>>:). Set LastRC for SFS:, SFU:, S3S:, and S3U:. Upgrade Status to Release Candidate 1

02/27/2021 - Beta 41(RC1) - Set LastRC for /GET:, GET:, and /GXR:

03/02/2021 - Beta 42(RC1) - Check if a port is open on a remote machine - TCP:RemoteHost:Port or UDP:RemoteHost:Port. IMPORTANT NOTE: UDP is connectionless and unreliable - I have included this functionality, but it cannot be trusted.  Use with Caution and Caveat.

03/05/2021 - Beta 43(RC1) - Close Ini File at the end of processing. Add LogHndl.Sync() after SAY: to control/force Log file Flushing better. Improve Unzip messages.

03/18/2021 - Beta 50(RC1) - Convert to Go1.16 (REQUIRED TO COMPILE THIS VERSION). Convert from AChoirX custom embedder to native GoLang Embed. Convert from GOPATH to Module. Improve UnZip Routine.

03/19/2021 - Beta 51(RC1) - Implement Syslog RFC3164 Format

03/19/2021 - Beta 52(RC1) - Add Syslog Type (SET:SyslogT=) of UDP or TCP. Improve Syslog Message format

03/24/2021 - Beta 53(RC1) - Improve Embedded Extraction Logic. Extract if AChoir.ACQ is not there. Allow other .ACQ files to be extracted and Run. Error Detection when Files Dissapear during processing

05/01/2021 - Beta 54(RC2) - Add File and Directory Delete Functions. DEL:<File To Delete> (Accepts WildCards) - Only Files in Subdirectories (Off of The AChoirX Root). CLN:<AChoirX Sub-Directory to Clean and Delete> - Only Subdirectories (Off of The AChoirX Root). This is to prevent accidental Deletion of files not related to the acquisition or toolkit

05/20/2021 - Beta 55(RC2) - Attempting to fix occasional Hang on Threads in the Wait Chain. The problem only happens on many small files. It may be related to deferring the Close.  Added Counters to troubleshoot the issue. Make Console Message and Log Levels the same.                     

05/21/2021 - Beta 56(RC2) - Separate the Debugging Counters to Isolate Better. Add Debug command Line Option - /DBG:<min>, <std>, <max>, <debug>

05/28/2021 - Beta 57(RC2) - Add Context and Timeout to AWS S3 upload for Upload hangs. Add Rudimentary Zip Routine - Must Use &FOR and cannot add to Zip.

05/30/2021 - Beta 58(RC2) - Expand and Improve Zip Routine: Allow Multiple Additions, Change Output Zip File Naming routines, and Add WildCards. 

08/15/2021 - Beta 59(RC2) - Small bug fix for determining current Disk Available (&DSA) if the Drive is not C: (Windows Only)

09/12/2021 - Beta 90(RC3) - Small bug fix for Delims. Add REX: Load Regular Expression Table. Add HST: Load Hash Table. Add Regular Expression Searching to CPS: (Copy by Signature). Add Hash Searching to CPS: (Copy by Signature)

09/19/2021 - Beta 91(RC3) - Add /B64:<Base64SEncodedIniFileOfAChoirCommands> - Allows a Base64 Encoded string to create an Ini File - work like the PowerShell -enc Parameter 

11/25/2021 - Beta 92(RC3) - Add Upload Retry Count (Default is 3)

12/27/2021 - Beta 93(RC3) - Escape percent signs

05/15/2022 - Beta 94(RC3) - Add Echo command

09/22/2022 - Beta 95(RC3) - Add CPU Limit Throttling

10/21/2022 - AChoirX v10.00.96 - Improve CPU Limit Throttling

10/24/2022 - AChoirX v10.00.97 - Add Native Registry Extraction

02/12/2023 - AChoirX v10.00.98 - Check for Collisions - Multiple collections at the same time, Improve Syslog (remove CRLFs)

03/10/2023 - AChoirX v10.00.99 - Minor improvement to CPS: (it ignores case now)

AChoirX v10.01.00 - Release 1.0, Add /Nam: to Specify Directory Name

04/02/2023 - AChoirX v10.01.01 - Release 1.01. Add FLT:<Filter Files> = Filter &LST and &FOR based on a Filter File, Add SET:Filter= to control how the filter functions, None = Remove the Filter, Incl or Excl = Filter is used to Include or Exclude entries, Full or Part = Filter is full or partial match, (Example: SET:Filter=Incl,Part = Filter data that has Partial Matches)

04/15/2023 - AChoirX v10.01.02 - Release 1.02. Fix Zip Bug when no directory is specified

04/16/2023 - AChoirX v10.01.03 - Release 1.03. More improvements in Zip - Fix Subdirectory Indexing

06/18/2023 - AChoirX v10.01.10 - Release 1.10. Add remote Multi-Handler (Server & Client Modes)

06/24/2023 - AChoirX v10.01.11 - Release 1.11. Improvements in remote Multi-Handler

06/30/2023 - AChoirX v10.01.12 - Release 1.12. Change Behavior: If Console or CLI was invoked. Drop back into Interractive Mode after INI: Processing

07/02/2023 - AChoirX v10.01.13 - Release 1.13. Get input from Stdin or TCP Server 

07/04/2023 - AChoirX v10.01.14 - Release 1.14. Con:Last - Display last 10 Console Messages 

07/06/2023 - AChoirX v10.01.15 - Release 1.15. Improvements in Con:Last 

07/08/2023 - AChoirX v10.01.16 - Release 1.16. Default CLI Console output Redirect to TCP 

07/15/2023 - AChoirX v10.01.17 - Release 1.17. SFTP improvements (cross platform). Add /XTR and XTR: - extracts the embedded toolkit

07/23/2023 - AChoirX v10.01.18 - Release 1.18. Cosmetic changes, Improvements in the embedded Default Scripts

12/23/2023 - AChoirX v10.01.19 - Release 1.19. Change behavior of --exestdout and --exestderr to Append Mode, Linux collection script improvements

12/24/2023 - AChoirX v10.01.20 - Release 1.20. Improve SRV/CLI display if there are errors

01/01/2024 - AChoirX v10.01.20a - Release 1.20a. Standardize Win, Lin, and OSX default embedded scripts across all executables.

01/01/2024 - AChoirX v10.01.20b - Release 1.20b. Update all upload scripts to support OSX.

01/07/2024 - AChoirX v10.01.20c - Release 1.20c. Add MacOS Support for ARM 64.

02/03/2024 - AChoirX v10.01.21 - Release 1.21. Occasionally the TCP STDOut file is not deleted (Add clear file to compensate)

03/18/2024 - AChoirX v10.01.21a - Release 1.21a. Include Beta Android Version. For more info: http://www.musectech.com/2024/03/achoirx-and-android-another-rabbit-hole.html

04/06/2024 - AChoirX v10.01.22 - Release 1.22 Dont allow Files to be zipped into themselves

04/20/2024 - AChoirX v10.01.23 - Release 1.23 Improvements in Registry Key Extraction

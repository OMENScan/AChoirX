# AChoirX
AChoirX is a port of AChoir in Go(lang) for Cross Platform forensic collection, processing, and analysis.

The goal of AChoirX is to create a simple scripting language for Targeted Collection (Live Response) on Windows, Linux, OSX (MacOS), and Android.  AChoirX (like AChoir) will also support scripting of extraction, parsing, and analysis of mounted forensic images.

The AChoirX installation contains a single AChoir collection script, which can run on all three primary versions of AChoirX (Windows, Linux, MacOS). This shows the power of the idea: A single consistent method to gather artifacts and telemetry from all three platforms using a single script.  AChoirX does not to abstract these artifacts - but rather collects the actual platform specific raw artifacts in a consistent manner using the same utility on all three platforms.  Please Note: While AChoirX also runs on Android - it is not yet supported in the default script.

AChoir scripts can be modified in nearly endless ways to meet the collection needs of different types of investigations. AChoir scripts, external programs, and other files (toolkits) can be located on the local drive, embedded into the AChoirX executables at compile time or appended to the AChoirX executables using the /PKR: option. Giving Incident Responders multiple, flexible ways to create custom collectors.

AChoirX can also run in interactive console mode, or across TCP using the built in Multi-Handler.  AChoirX also natively supports common upload methods such as SFTP or S3.  Since AChoirX runs on Windows, Linux, MacOS, and Android - All of these functions, including the multi-handler (client and server modes) and the upload methods, work regardless of platform.

# Installing / Running
You can install or just run AChoirX in multiple ways.  AChoirX is designed to be flexible to accommodate multiple levels of default or custom collectors. Some of the options are:

1. Download and run the <b>AChoirX-inst.exe</b> installer from this repo (<b>BE SURE TO RUN THIS AS ADMIN</b>)
   - Running the installer will load AChoirX onto your machine (basically just unzipping the AChoirX components onto your drive).  This install contains the source code, some FOSS utilities, and lots of AChoirX scripts.  You can simply use the default embedded script in the Windows, MacOS, or Linux AChoirX executables, or call any of the other embedded scripts using the <b>/INI:</b> option. The install contains ready to run Windows, Linux, and MacOS collectors.
   - The installer will then run the AChoirX builder script (<b>Build.ACQ</b>) which downloads several additional FOSS Utilities, packages them into a (.zip) toolkit and then creates an extensive Windows collector by embedding the toolkit. This new collector (<b>AChWinAll.exe</b>) uses both native AChoirX capabilities and the additional downloaded FOSS tools to extend the collector's capabilities. You can do the same thing by using the (<b>/PKR:</b>) option of AChoirX.
2. Download and run the default Windows, Linux, and MacOS AChoirX collectors from this repo as-is.  This is the fastest way to get started with AChoirX - Just download and run them.  The executables are self contained, No assembly required.
3. Download and run (<b>AS ADMIN</b>) the default AChoirX collector from this repo with the <b>/BLD</b> option to run the AChoirX builder script (<b>Build.ACQ</b>) which downloads several additional FOSS Utilities, Packages them into a (.zip) toolkit and then creates an extensive Windows collector by embedding the toolkit. This new collector (<b>AChWinAll.exe</b>) uses both native AChoirX capabilities and the additional downloaded FOSS tools to extend the collector's capabilities. You can do the same thing by using the (<b>/PKR:</b>) option of AChoirX.
4. Create your own AChoirX (.zip) toolkit and embed it directly into AChoirX using the <b>/PKR:</b> option.  This is an easy way to create your own custom collectors.  Take a look at the <b>ToolsBuiltWithAChoirX</b> section of this repo for some ideas of the types of collectors you can create.
5. Compile AChoirX using GO, and including the default (or your own) toolkit to create a custom collector. AChoirX is Free and Open Source Software.

Choose your own adventure!

# Disclaimer
This repository represents a Project and not a Product.

This software is furnished "as-is". We provide no warranty whatsoever, whether express, implied, or statutory, including, but not limited to, any warranty of merchantability or fitness for a particular reason or purpose, or any warranty that the software will be error-free.

In no respect shall the author or distributors of AChoirX incur any liability for any damages, including, but not limited to, direct, indirect, special, or consequential damages arising out of, resulting from, or any way connected to the use of this software, whether or not based upon warranty, contract, tort, or otherwise; whether or not injury was sustained by persons or property or otherwise; and whether or not loss was sustained from, or arose out of, the results of, the software, or any services that may be provided by the software.

You agree to install and use this software solely at your own risk.

AChoirX is an Open Source Project and carries no formal support expressed or implied.


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

12/27/2020 - Beta 36. Change Conditional Logic to only count a single occurrence of &FOR and &LST comparisons. This prevents the need for multiple END: statements.
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

06/30/2023 - AChoirX v10.01.12 - Release 1.12. Change Behavior: If Console or CLI was invoked. Drop back into Interactive Mode after INI: Processing

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

04/22/2024 - AChoirX v10.01.24 - Release 1.24 Fix CopyPath=Full Bug 

05/19/2024 - AChoirX v10.01.50 - Release 1.50 No code changes per se. Refactoring of the litany of Scripts - Consolidating wherever it makes sense. 

06/03/2024 - AChoirX v10.01.51 - Release 1.51 Change behavior of &LST when it is used with FOR: (i.e. to do FOR: on a list of directories). When in Looping mode, Append &FOR file names into ForFiles instead of overwriting them 

06/20/2024 - AChoirX v10.01.52 - Release 1.52 - Change behavior of &LST and &FOR on parsing error - &LST and &FOR will work even on parse error. Add Set:ParseQuote=Strict or Lazy

06/20/2024 - AChoirX v10.01.53 - Release 1.53 - Change REG: output file name to replace any invalid chars with "-"

06/20/2024 - AChoirX v10.01.54 - Release 1.54 - Add LST: and FOR: Counters

08/17/2024 - AChoirX v10.01.55 - Release 1.55 - Add INC: (Include an INI - Allowing Nested INI Files)

09/03/2024 - AChoirX v10.01.56 - Release 1.56 - Add OPTIONAL output file name to REG: to allow all extractions to go to the same CSV - Close Registry Key Properly so it can be unloaded

09/29/2024 - AChoirX v10.01.57 - Release 1.57 - Add NCP: NTFS Raw Copy - Only implemented in Windows - Not Applicable to Linux, MacOS, or Android - Most of the code for this was copied from: https://github.com/kmahyyg/go-rawcopy

10/11/2024 - AChoirX v10.01.58 - Release 1.58 - Change Time Display to UTC. Fix edge case when Regexp expands "$" like with $MFT or $Logfile, etc...

11/10/2024 - AChoirX v10.01.58a - Release 1.58a - Componentize, Combine, and Consolidate Scripts

11/30/2024 - AChoirX v10.01.59 - Release 1.59 - Change XIT: to run non-Blocked (EXA:) - This change in behavior is designed to allow AChCleanr to run after AChoirX exits. Add &MyE - For This Program's Executable name. Add experimental AChCleanr - Which dissolves most of AChoirX after it runs.  The forensic collection data will remain, but most other components will be erased.  This is to prevent a Bad Actor from using AChoirX left behind on the drive as a LOLBin.

10/24/2025 - AChoirX v10.01.60 - Release 1.60 - Prevent Globbing errors from exiting walking directories 

11/04/2025 - AChoirX v10.01.61 - Release 1.61 - Add Regex to Filter behavior (Filter=Regx)

11/14/2025 - AChoirX v10.01.62 - Release 1.62 - Minor change in Zipping error

11/27/2025 - AChoirX v10.01.75 - Release 1.75 - Add a Packager and UnPacker to allow creating customized versions                      of AChoirX by Appending a Zip File to the end of the program instead of requiring a recompile with an embedded zip.  The logic is: 1. First - Try local Scripts and program, 2. Next - Try to unzip any packaged zip, 3. Last - Unzip the embedded zip compiled into the executable. This creates an extremely flexible way to deploy custom programs and scripts.

12/01/2025 - AChoirX v10.01.76 - Release 1.76 - Minor Log Copy change. 

12/13/2025 - AChoirX v10.01.77 - Release 1.77 - Do not allow multiple Toolkits to be embedded. A security feature to help prevent abuse of AChoirX toolkit embedding.

12/14/2025 - AChoirX v10.01.78 - Release 1.78 - Preserve relative paths when copying/uploading files using double-globbing.

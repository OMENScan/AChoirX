# AChoirX
 AChoirX is a port of AChoir in Go(lang) for Cross Platform forensic collection, processing, and analysis.

# The Plan:
1. Port all functions from AChoir to AChoirX that can be cross-platform supported to GO.  This will give a basic set of functionality, but WILL NOT have all Windows specific functionality of the C version of AChoir.  I recommend that the C Version on windows be used for the foreseeable future.  I will continue to Update and enhance the Windows C Version.  For now, the C version is the flagship version of AChoir for Windows.
2. Add the functions that are not available using Standard Go, but are still cross-platform.  I will rely on third party libraries, or create my own libraries to ensure that functionality not available in native Go (but important for all platforms) will be implemented in AChoirX.
3. Add platform specific code to add functionality that is either not available in all platforms, or is handled differently in different platforms. AChoirX will contain shared code, and platform specific code.  I'll try to make this as transparent as possible, but I anticipate that some things will simply remain platform specific.

The goal of AChoirX is to create a simple scripting language for Targeted Collection (Live Response) on Windows, Linux, and OSX.  AChoirX (like AChoir) will also support scripting of extraction, parsing, and analysis of mounted forensic images.

In the latest Alpha release I have created a single AChoir collection script, which can run on all three versions of AChoir. This shows the power of the idea: A single consistent method to gather artifacts and telemetry from all three platforms using a single script.  AChoirX does not to abstract these artifacts - but rather collect the actual platform specific raw artifacts in a consistent manner using the same utility on all three platforms.

I hope that eventually this idea will help forensicators move across multiple platforms without the need for separate tools requiring separate syntax.

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

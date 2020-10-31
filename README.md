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
  
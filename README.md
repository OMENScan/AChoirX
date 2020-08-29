# AChoirX
 ReWrite of AChoir in Go for Cross Platform

# The Plan:
1. Convert all functions that can be cross-platform supported to GO.  this will give a basic set of functionality, but WILL NOT have all functionality of the C version of AChoir.  I recommend that the C Version be used for the forseeable future.  I will continue to Update and enhance the C Version.  For now, the C version is the flagship version of AChoir for Windows.
2. Add the functions that are not available using Standard Go, but are still cross-platform
3. Add platform specific code to addfunctionality that is either not available in all platforms, or is handled differenty in different platforms.

The goal of AChoirX is to create a simple scripting language for Targeted Collection (Live Response) on Windows, Linux, and OSX.  In theory - Single AChoir collection scripts can be created that can run on all three versions of AChoir. The goal is to allow a single consistent method to gather artifacts and telemetry from all three platforms.  My goal is not to abstract these artifacts - but rather to collect the actual raw artifacts in a consistent manner using the same utility on all three platforms.

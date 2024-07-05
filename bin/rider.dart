import 'dart:io';

import 'package:package_config/package_config.dart';
import 'package:path/path.dart' as path;

void main(List<String> arguments) async {
  // Locate the package configuration file
  final pwd = Directory.current.path;
  print('Current directory: $pwd/.dart_tool/package_config.json');
  final packageConfigFile = File('$pwd/.dart_tool/package_config.json');
  if (!packageConfigFile.existsSync()) {
    print(
      'Package configuration not found. Ensure you are running this within a Dart/Flutter project.',
    );
    exit(1);
  }

  // Load the package configuration
  final packageConfig = await loadPackageConfig(packageConfigFile);
  final package = packageConfig['rider'];
  if (package == null) {
    print('Package "rider" not found in package configuration.');
    exit(1);
  }

  // Determine the executable path based on the platform
  final packageRoot = package.root.toFilePath();
  String executablePath;

  if (Platform.isMacOS) {
    executablePath =
        path.join(packageRoot, 'lib', 'src', 'cli', 'bin', 'darwin', 'rider');
  } else if (Platform.isLinux) {
    executablePath =
        path.join(packageRoot, 'lib','src', 'cli', 'bin', 'linux', 'rider');
  } else if (Platform.isWindows) {
    executablePath =
        path.join(packageRoot, 'lib','src', 'cli', 'bin', 'windows', 'rider.exe');
  } else {
    print('Unsupported platform');
    exit(1);
  }

  // If no argument is given use init as default
  if (arguments.isEmpty) {
    arguments = List.from(["init"]);
  }

  // Ensure the executable has execution permission on Unix-like systems
  if (!Platform.isWindows) {
    await Process.run('chmod', ['+x', executablePath]);
  }

  // Execute the binary with the given arguments
  final result = await Process.start(
    executablePath,
    arguments,
    mode: ProcessStartMode.inheritStdio,
  );

  // Exit with the binary's exit code
  exit(await result.exitCode);
}

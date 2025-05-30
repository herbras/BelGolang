#!/usr/bin/env node
import path from 'path';
import fs from 'fs';
import os from 'os';
import { fileURLToPath } from 'url';

// ES module compatibility
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Runtime detection
const isBun = typeof Bun !== 'undefined' && typeof Bun.spawn === 'function';
const isDeno = typeof Deno !== 'undefined' && typeof Deno.Command === 'function';
const isNode = typeof process !== 'undefined' && !!process.versions.node;

function getPlatform() {
  const platform = os.platform();
  const arch = os.arch();
  
  const platformMap = {
    'darwin': 'darwin',
    'linux': 'linux', 
    'win32': 'win32'
  };
  
  const archMap = {
    'x64': 'x64',
    'arm64': 'arm64',
    'aarch64': 'arm64'
  };
  
  return {
    platform: platformMap[platform] || platform,
    arch: archMap[arch] || 'x64'
  };
}

function getBinaryPath() {
  const { platform, arch } = getPlatform();
  const binaryName = platform === 'win32' ? 'salat.exe' : 'salat';
  const binaryPath = path.join(__dirname, 'bin', `${platform}-${arch}`, binaryName);
  
  if (!fs.existsSync(binaryPath)) {
    console.error(`❌ Binary not found for ${platform}-${arch}`);
    console.error(`Expected: ${binaryPath}`);
    console.error('\nSupported platforms:');
    console.error('  - darwin-x64, darwin-arm64');
    console.error('  - linux-x64, linux-arm64');
    console.error('  - win32-x64, win32-arm64');
    process.exit(1);
  }
  
  return binaryPath;
}

// Unified spawn wrapper for all runtimes
async function runSalatCLI(binaryPath, args) {
  if (isNode) {
    const { spawn } = await import('child_process');
    const child = spawn(binaryPath, args, {
      stdio: 'inherit',
      windowsHide: false
    });
    
    child.on('error', (error) => {
      console.error(`❌ Failed to start salat CLI: ${error.message}`);
      process.exit(1);
    });
    
    child.on('exit', (code) => {
      process.exit(code || 0);
    });
    
    return child;
  }
  
  if (isBun) {
    try {
      const proc = Bun.spawn([binaryPath, ...args], { 
        stdio: ['inherit', 'inherit', 'inherit'] 
      });
      const exitCode = await proc.exited;
      process.exit(exitCode);
    } catch (error) {
      console.error(`❌ Failed to start salat CLI: ${error.message}`);
      process.exit(1);
    }
  }
  
  if (isDeno) {
    try {
      const cmd = new Deno.Command(binaryPath, { 
        args: args,
        stdin: 'inherit',
        stdout: 'inherit',
        stderr: 'inherit'
      });
      const child = cmd.spawn();
      const status = await child.status;
      Deno.exit(status.code);
    } catch (error) {
      console.error(`❌ Failed to start salat CLI: ${error.message}`);
      Deno.exit(1);
    }
  }
  
  console.error('❌ Unsupported JavaScript runtime');
  console.error('Supported runtimes: Node.js, Bun, Deno');
  (typeof Deno !== 'undefined' ? Deno.exit : process.exit)(1);
}

async function main() {
  const args = (typeof Deno !== 'undefined' ? Deno.args : process.argv.slice(2)) || [];
  
  if (args.includes('--verify-install')) {
    const binaryPath = getBinaryPath();
    const runtime = isBun ? 'Bun' : isDeno ? 'Deno' : 'Node.js';
    console.log(`✅ Salat CLI installed successfully for ${getPlatform().platform}-${getPlatform().arch}`);
    console.log(`Runtime: ${runtime}`);
    console.log(`Binary location: ${binaryPath}`);
    console.log('\nRun "salat setup" to get started!');
    return;
  }

  const binaryPath = getBinaryPath();
  await runSalatCLI(binaryPath, args);
}

// ES module entry point check - robust cross-platform and multi-Node.js compatible
const isMainModule = async () => {
  try {
    // Method 1: Direct comparison (works for most cases)
    if (import.meta.url === `file://${process.argv[1]}`) {
      return true;
    }
    
    // Method 2: EndsWith check (fallback)
    if (import.meta.url.endsWith(process.argv[1])) {
      return true;
    }
    
    // Method 3: Path normalization for cross-platform compatibility
    const { fileURLToPath } = await import('url');
    const metaPath = fileURLToPath(import.meta.url);
    const argvPath = process.argv[1];
    
    // Direct path comparison
    if (metaPath === argvPath) {
      return true;
    }
    
    // Method 4: Normalize both paths and compare (handles different path separators)
    const path = await import('path');
    const normalizedMeta = path.resolve(metaPath);
    const normalizedArgv = path.resolve(argvPath);
    
    if (normalizedMeta === normalizedArgv) {
      return true;
    }
    
    // Method 5: Compare just the filename and parent directory (handles different Node.js installations)
    const metaBasename = path.basename(metaPath);
    const argvBasename = path.basename(argvPath);
    const metaParent = path.basename(path.dirname(metaPath));
    const argvParent = path.basename(path.dirname(argvPath));
    
    // If both are index.js in salat-cli directory, assume it's the same file
    if (metaBasename === 'index.js' && argvBasename === 'index.js' && 
        metaParent === 'salat-cli' && argvParent === 'salat-cli') {
      return true;
    }
    
    return false;
  } catch (error) {
    // Fallback: if all else fails, assume we're being run directly
    console.warn('Entry point detection failed, assuming direct execution');
    return true;
  }
};

// Use async wrapper for isMainModule
(async () => {
  if (await isMainModule()) {
    main().catch(error => {
      console.error('❌ Unexpected error:', error.message);
      (typeof Deno !== 'undefined' ? Deno.exit : process.exit)(1);
    });
  }
})();

export { getBinaryPath, getPlatform, runSalatCLI };
#!/usr/bin/env node
const { execFileSync } = require('child_process');
const ffi = require('ffi-napi');
const path = require('path');
const os = require('os');

const getLibrary = () => {
  const platform = os.platform();

  console.log(path.join(__dirname, 'bin', 'git-manager.dylib'))
  const libPath = {
    linux: path.join(__dirname, 'bin','git-manager.so'),
    darwin: path.join(__dirname, 'bin', 'git-manager.dylib'),
    win32: path.join(__dirname, 'bin','git-manager.dll'),
  }[platform];

  if (!libPath) throw new Error('Platform not supported');
  return libPath;
};

const args = process.argv.slice(2);
execFileSync(getLibrary(), args, { stdio: 'inherit' });

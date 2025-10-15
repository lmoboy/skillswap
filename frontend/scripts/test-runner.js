#!/usr/bin/env node

const { execSync } = require('child_process');
const path = require('path');

// Colors for console output
const colors = {
  reset: '\x1b[0m',
  bright: '\x1b[1m',
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  magenta: '\x1b[35m',
  cyan: '\x1b[36m'
};

function log(message, color = 'reset') {
  console.log(`${colors[color]}${message}${colors.reset}`);
}

function runCommand(command, description) {
  log(`\n${colors.cyan}Running: ${description}${colors.reset}`);
  log(`${colors.yellow}Command: ${command}${colors.reset}\n`);
  
  try {
    execSync(command, { 
      stdio: 'inherit',
      cwd: process.cwd(),
      shell: true
    });
    log(`${colors.green}✓ ${description} completed successfully${colors.reset}`);
    return true;
  } catch (error) {
    log(`${colors.red}✗ ${description} failed${colors.reset}`);
    log(`${colors.red}Error: ${error.message}${colors.reset}`);
    return false;
  }
}

function main() {
  const args = process.argv.slice(2);
  const testType = args[0] || 'all';
  
  log(`${colors.bright}${colors.blue}🧪 Frontend Test Runner${colors.reset}`);
  log(`${colors.cyan}Running tests: ${testType}${colors.reset}\n`);
  
  const startTime = Date.now();
  let allPassed = true;
  
  switch (testType) {
    case 'unit':
      allPassed = runCommand('bun run test:unit', 'Unit Tests');
      break;
      
    case 'integration':
      allPassed = runCommand('bun run test:integration', 'Integration Tests');
      break;
      
    case 'e2e':
      allPassed = runCommand('bun run test:e2e', 'End-to-End Tests');
      break;
      
    case 'coverage':
      allPassed = runCommand('bun run test:coverage', 'Test Coverage');
      break;
      
    case 'all':
    default:
      log(`${colors.magenta}Running all test suites...${colors.reset}`);
      
      const unitPassed = runCommand('bun run test:unit', 'Unit Tests');
      const integrationPassed = runCommand('bun run test:integration', 'Integration Tests');
      const e2ePassed = runCommand('bun run test:e2e', 'End-to-End Tests');
      
      allPassed = unitPassed && integrationPassed && e2ePassed;
      break;
  }
  
  const endTime = Date.now();
  const duration = ((endTime - startTime) / 1000).toFixed(2);
  
  log(`\n${colors.bright}${colors.blue}Test Summary${colors.reset}`);
  log(`${colors.cyan}Duration: ${duration}s${colors.reset}`);
  
  if (allPassed) {
    log(`${colors.green}${colors.bright}🎉 All tests passed!${colors.reset}`);
    process.exit(0);
  } else {
    log(`${colors.red}${colors.bright}❌ Some tests failed${colors.reset}`);
    process.exit(1);
  }
}

// Handle process signals
process.on('SIGINT', () => {
  log(`\n${colors.yellow}Test run interrupted${colors.reset}`);
  process.exit(1);
});

process.on('SIGTERM', () => {
  log(`\n${colors.yellow}Test run terminated${colors.reset}`);
  process.exit(1);
});

main();

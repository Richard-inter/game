#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const { spawn } = require('child_process');

// Service configuration
const services = [
    { name: 'Game Service', port: 9090, type: 'gRPC', file: 'game-service' },
    { name: 'API Service', port: 8080, type: 'HTTP', file: 'api-service' },
    { name: 'WebSocket Service', port: 8081, type: 'WebSocket', file: 'websocket-service' },
    { name: 'TCP Service', port: 8082, type: 'TCP', file: 'tcp-service' }
];

// Check if service is running
function checkService(service) {
    return new Promise((resolve) => {
        const net = require('net');
        const socket = new net.Socket();
        
        const timeout = setTimeout(() => {
            socket.destroy();
            resolve({ ...service, status: 'stopped', pid: null });
        }, 1000);

        socket.connect(service.port, 'localhost', () => {
            clearTimeout(timeout);
            socket.destroy();
            resolve({ ...service, status: 'running', pid: null });
        });

        socket.on('error', () => {
            clearTimeout(timeout);
            resolve({ ...service, status: 'stopped', pid: null });
        });
    });
}

// Get process ID for a service
function getServicePID(serviceFile) {
    return new Promise((resolve) => {
        const ps = spawn('pgrep', ['-f', `./cmd/${serviceFile}`]);
        let output = '';
        
        ps.stdout.on('data', (data) => {
            output += data.toString();
        });
        
        ps.on('close', (code) => {
            const pids = output.trim().split('\n').filter(pid => pid.trim());
            resolve(pids.length > 0 ? parseInt(pids[0]) : null);
        });
    });
}

// Main monitoring function
async function monitorServices() {
    console.log('🔍 Checking service status...\n');
    
    let runningCount = 0;
    
    for (const service of services) {
        const serviceStatus = await checkService(service);
        const pid = await getServicePID(service.file);
        serviceStatus.pid = pid;
        
        if (serviceStatus.status === 'running') {
            runningCount++;
            console.log(`✅ ${serviceStatus.name} (${serviceStatus.type}) - Port ${serviceStatus.port} - PID: ${pid || 'Unknown'}`);
        } else {
            console.log(`❌ ${serviceStatus.name} (${serviceStatus.type}) - Port ${serviceStatus.port} - Stopped`);
        }
    }
    
    console.log(`\n📊 Summary: ${runningCount}/${services.length} services running\n`);
    
    // Show quick commands
    console.log('🚀 Quick Commands:');
    console.log('   Start all: make docker-up-infra && use VS Code launch configs');
    console.log('   Check ports: lsof -i :8080,8081,8082,9090');
    console.log('   Kill all: pkill -f "game-service|api-service|websocket-service|tcp-service"');
    
    return runningCount === services.length;
}

// Export for use in other scripts
module.exports = { monitorServices, services };

// Run if called directly
if (require.main === module) {
    monitorServices().then(allRunning => {
        process.exit(allRunning ? 0 : 1);
    });
}

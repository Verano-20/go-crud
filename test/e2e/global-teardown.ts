import { FullConfig } from '@playwright/test';

async function globalTeardown(config: FullConfig): Promise<void> {
  console.log('🧹 Starting global teardown for E2E tests...');
  
  // Add any cleanup logic here if needed
  // For example, cleaning up test data, stopping services, etc.
  
  console.log('✅ Global teardown completed');
}

export default globalTeardown;

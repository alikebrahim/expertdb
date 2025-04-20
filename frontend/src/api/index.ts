import * as auth from './auth';
import * as experts from './experts';
import * as requests from './requests';
import * as documents from './documents';
import * as engagements from './engagements';
import * as phases from './phases';
import * as statistics from './statistics';
import * as areas from './areas';
import * as users from './users';
import * as backup from './backup';
import { apiClient, createApiClient } from './client';

export {
  apiClient,
  createApiClient,
};

export default {
  auth,
  experts,
  requests,
  documents,
  engagements,
  phases,
  statistics,
  areas,
  users,
  backup,
};
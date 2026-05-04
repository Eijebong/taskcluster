import { describe, it, before, after, beforeEach, afterEach } from 'node:test';
import assert from 'assert';
import { APIBuilder } from '../src/index.js';
import testing from '@taskcluster/lib-testing';

describe(testing.suiteName(), function() {
  it('middleware is exported', function() {
    assert(APIBuilder.middleware);
  });
});

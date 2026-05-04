import { describe, it, before, after, beforeEach, afterEach } from 'node:test';
import { cleanRouteAndParams } from '../src/utils.js';
import testing from '@taskcluster/lib-testing';
import { strict as assert } from 'assert';

describe(testing.suiteName(), function() {
  describe('cleanRouteAndParams', function() {
    it('for a plain route', function() {
      assert.deepEqual(
        cleanRouteAndParams('/ab/cd'),
        ['/ab/cd', [], []]);
    });

    it('for a route with "regular" params', function() {
      assert.deepEqual(
        cleanRouteAndParams('/foo/:foo/bar/:bar'),
        ['/foo/<foo>/bar/<bar>', ['foo', 'bar'], []]);
    });

    it('for a route with an optional param', function() {
      assert.deepEqual(
        cleanRouteAndParams('/foo/:foo/bar/:bar?'),
        ['/foo/<foo>/bar/<bar>', ['foo', 'bar'], ['bar']]);
    });
  });
});

// jshint ignore: start
/// <reference path="../../types/realm.d.ts" />

'use strict';

import { RSFunctionQueryData } from '../../types/types';
import { ReactiveSearch } from '../../index';

// @ts-ignore
exports = async (payload: any) => {
	// @ts-expect-error
	const query: RSFunctionQueryData = EJSON.parse(payload.body.text());
	const { config, searchQuery } = query;
	const client = context.services.get('mongodb-atlas');

	const reactiveSearch = new ReactiveSearch({
		client,
		database: config.database,
		documentCollection: config.documentCollection,
	});

	const results = await reactiveSearch.query(searchQuery);
	return results;
};

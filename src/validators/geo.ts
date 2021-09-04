import { GeoValue } from 'src/types/types';

export const validateGeoValue = (val: GeoValue) => {
	if (!val.location && !val.geoBoundingBox) {
		throw new Error(`Invalid object`);
	}
};

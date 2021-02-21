#include <stdint.h>
#include <string.h>
#include "hashmap.c"

#define IS_POINTER(x) (((uint64_t)(x) & 1) == 0)

struct solestr {
	char* value;
};

static char * solestring_value_pack(char *v) {
	size_t len = strlen(v);
#ifdef SOLESTRING_PACK
	// Pointer
	if (len > sizeof(char *)-2) { // one byte for tag and one for \0
#endif
		char *np = malloc((len+1)*sizeof(char));
		memcpy(np, v, len+1);
		return np;
#ifdef SOLESTRING_PACK
	}
	// Tagged pointer
	char *tp;
	char *p = (char *)&tp; // reinterpret tag pointer itself as byte array
	*p = 1; // tag in the first byte (assuming little endian)
	memcpy(p+1, v, len+1);
	return tp;
#endif
}

static char * solestring_value_unpack(char **vp) {
	return IS_POINTER(*vp) ? *vp : (char *)vp+1;
}

static int solestr_compare(const void *a, const void *b, void *udata) {
	struct solestr *sa = (struct solestr *)a;
	struct solestr *sb = (struct solestr *)b;
	return strcmp(solestring_value_unpack(&sa->value), solestring_value_unpack(&sb->value));
}

static uint64_t solestr_hash(const void *item, uint64_t seed0, uint64_t seed1) {
	struct solestr *ss = (struct solestr *)item;
	char *v = solestring_value_unpack(&ss->value);
	return hashmap_sip(v, strlen(v), seed0, seed1);
}

static bool solestr_iter_free(const void *item, void *udata) {
	struct solestr *ss = (struct solestr *)item;
	struct hashmap *hmap = udata;
	if (IS_POINTER(ss->value)) {
		free(ss->value);
	}
	hashmap_delete(hmap, ss);
	return true;
}

struct hashmap * hmap_new() {
	return hashmap_new(sizeof(struct solestr), 0, 0, 0,
					   solestr_hash, solestr_compare, NULL);
}

char * hmap_get(struct hashmap *hmap, char *s) {
	struct solestr *ss = hashmap_get(hmap, &(struct solestr){ .value=s });
	return ss ? ss->value : NULL;
}

bool hmap_put(struct hashmap *hmap, char *s) {
	struct solestr *ss = hashmap_set(hmap, &(struct solestr){ .value=solestring_value_pack(s) });
	return ss || !hashmap_oom(hmap);
}

void hmap_free(struct hashmap *hmap) {
	hashmap_scan(hmap, solestr_iter_free, hmap);
	hashmap_free(hmap);
}

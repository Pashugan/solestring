#include <stdint.h>
#include <string.h>
#include "hashmap.c"

#define IS_POINTER(x) (((uint64_t)(x) & 1) == 0)

struct solestr {
	char* value;
};

static char * solestring_value_pack(char *v) {
	size_t len = strlen(v);
	if (len > 6) {
		char *np = (char *)malloc((len+1)*sizeof(char));
		strncpy(np, v, len+1);
		return np;
	}
	char *tp; // tagged pointer
	char *p = (char *)&tp;
	*p = 1;
	strncpy(p+1, v, len+1);
	return tp;
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

#include <stdint.h>
#include <string.h>
#include "hashmap.c"

#define IS_POINTER(x) (((uint64_t)(x) & 1) == 0)

struct solestr {
	char* value;
};

static void solestring_value_pack(char **vp) {
	size_t len = strlen(*vp);
	if (len > 6) {
		return;
	}
	uint64_t i = 1;
	char *ip = (char *)&i;
	strncpy(ip+1, *vp, len);
	*(ip+1+len) = '\0';
	*vp = (char *)i;
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
	return ss ? solestring_value_unpack(&ss->value) : NULL;
}

bool hmap_put(struct hashmap *hmap, char *s) {
	solestring_value_pack(&s);
	struct solestr *ss = hashmap_set(hmap, &(struct solestr){ .value=s });
	return ss || !hashmap_oom(hmap);
}

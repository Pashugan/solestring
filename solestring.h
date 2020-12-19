#include <string.h>
#include "hashmap.c"

struct solestr {
    char* value;
};

static int solestr_compare(const void *a, const void *b, void *udata) {
    return strcmp(((const struct solestr *)a)->value, ((const struct solestr *)b)->value);
}

static uint64_t solestr_hash(const void *item, uint64_t seed0, uint64_t seed1) {
    const struct solestr *ss = item;
    return hashmap_sip(ss->value, strlen(ss->value), seed0, seed1);
}

struct hashmap * hmap_new() {
    return hashmap_new(sizeof(struct solestr), 0, 0, 0,
                       solestr_hash, solestr_compare, NULL);
}

char *hmap_get(struct hashmap *hmap, char *s) {
    struct solestr *ss = hashmap_get(hmap, &(struct solestr){ .value=s });
    return ss ? ss->value : NULL;
}

bool hmap_put(struct hashmap *hmap, char *s) {
    struct solestr *ss = hashmap_set(hmap, &(struct solestr){ .value=s });
    return ss || !hashmap_oom(hmap);
}

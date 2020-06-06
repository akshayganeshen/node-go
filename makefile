all: lib mod
clean: clean-lib clean-mod clean-gyp

SOURCE_LIB_FILES = \
  $(wildcard lib/*.go) \
  $(wildcard lib/**/*.go) \
  $(wildcard lib/**/*.c) \
  $(wildcard lib/include/*.h)
SOURCE_MOD_FILES = \
  $(wildcard src/*.cpp) \
  $(wildcard src/*.h)

SOURCE_GYP_FILES = binding.gyp

TARGET_LIB_HEADER = lib/libgo.h
TARGET_LIB_ARCHIVE = lib/libgo.a

TARGET_GYP_BUILDDIR = build
TARGET_GYP_MODULE = nodego
TARGET_GYP_TYPE = Release
TARGET_GYP_STAMP = gyp.config.stamp

TARGET_GYP_TYPEDIR = $(TARGET_GYP_BUILDDIR)/$(TARGET_GYP_TYPE)
TARGET_GYP_MODFILE = $(TARGET_GYP_MODULE).node
TARGET_GYP_OUTPUT = $(TARGET_GYP_TYPEDIR)/$(TARGET_GYP_MODFILE)

lib: $(TARGET_LIB_HEADER) $(TARGET_LIB_ARCHIVE)
mod: $(TARGET_GYP_OUTPUT)

$(TARGET_LIB_HEADER): $(TARGET_LIB_ARCHIVE)
$(TARGET_LIB_ARCHIVE): $(SOURCE_LIB_FILES)
	go build -buildmode=c-archive -o "$@" "./lib"

$(TARGET_GYP_OUTPUT): $(TARGET_LIB_HEADER) $(TARGET_LIB_ARCHIVE)
$(TARGET_GYP_OUTPUT): $(SOURCE_MOD_FILES)
$(TARGET_GYP_OUTPUT): $(TARGET_GYP_STAMP)
	+MAKEFLAGS="$(MAKEFLAGS)" \
	  node-gyp build --make="$(MAKE)"

$(TARGET_GYP_STAMP): $(SOURCE_GYP_FILES)
	node-gyp configure
	touch "$@"

clean-lib:
	rm -f "$(TARGET_LIB_HEADER)" "$(TARGET_LIB_ARCHIVE)"

clean-mod:
	rm -rf "$(TARGET_GYP_BUILDDIR)"

clean-gyp:
	rm -f "$(TARGET_GYP_STAMP)"

.PHONY: all lib mod
.PHONY: clean clean-lib clean-mod clean-gyp

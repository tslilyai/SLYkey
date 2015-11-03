CC=gcc
CXX=g++
RM=rm -f
CPPFLAGS=""
LDFLAGS=""
LDLIBS=""

SRCS=user_api.cc test_user_api.cc
OBJS=$(subst .cc,.o,$(SRCS))

all: test_user_api 

test_user_api: $(OBJS)
    $(CXX) $(LDFLAGS) -o test_user_api $(OBJS) $(LDLIBS)

depend: .depend

.depend: $(SRCS)
    rm -f ./.depend
    $(CXX) $(CPPFLAGS) -MM $^>>./.depend;

clean:
	$(RM) $(OBJS)

dist-clean: clean
    $(RM) *~ .depend

include .depend

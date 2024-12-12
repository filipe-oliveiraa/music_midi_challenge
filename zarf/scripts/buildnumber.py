#!/usr/bin/env python3

# Computes a build number that is:
# D[DD]HH
# Where D is # days since our epoch date of 9/02/2024 UTC
# and HH is the hour of the day, prefixed with '0' if < 10
# e.g. if NOW is 5/28/2018 5:30am
#   => 305
import datetime

epoch = datetime.datetime(2024, 9, 2, 0, 0, 0)
d1 = datetime.datetime.now(datetime.UTC).replace(tzinfo=None)
delta = d1 - epoch
print("%d%02d" % (delta.days, d1.hour))

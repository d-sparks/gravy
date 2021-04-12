import math


def z_score_or_zero(x, mu, sigma):
    """
    Returns the z score unless sigma is prohibitively small, in which case
    returns 0.0. TODO: Move to a helper package.
    """
    if sigma == None or sigma < 1E-6:
        return 0.0
    return (x - mu) / sigma


def sqrt_or_zero(x):
    """
    Returns hte square root of x if x >= 0.0. TODO: Move to a helper package.
    """
    return math.sqrt(x) if x >= 0.0 else 0.0

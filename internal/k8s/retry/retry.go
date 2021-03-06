// Package retry is a very minimal wrapper around the "k8s.io/client-go/util/retry" package.
package retry

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	retry "k8s.io/client-go/util/retry"
)

// DefaultRetry is the recommended retry for a conflict where multiple clients
// are making changes to the same resource.
var DefaultRetry = retry.DefaultRetry

// DefaultBackoff is the recommended backoff for a conflict where a client
// may be attempting to make an unrelated modification to a resource under
// active management by one or more controllers.
var DefaultBackoff = retry.DefaultBackoff

// OnError allows the caller to retry fn in case the error returned by fn is retriable
// according to the provided function. backoff defines the maximum retries and the wait
// interval between two retries.
func OnError(backoff wait.Backoff, retriable func(error) bool, fn func() error) error {
	return retry.OnError(backoff, retriable, fn)
}

// OnConflict is used to make an update to a resource when you have to worry about
// conflicts caused by other code making unrelated updates to the resource at the same
// time. fn should fetch the resource to be modified, make appropriate changes to it, try
// to update it, and return (unmodified) the error from the update function. On a
// successful update, OnConflict will return nil. If the update function returns a
// "Conflict" error, OnConflict will wait some amount of time as described by
// backoff, and then try again. On a non-"Conflict" error, or if it retries too many times
// and gives up, OnConflict will return an error to the caller.
//
//     err := retry.OnConflict(retry.DefaultRetry, func() error {
//         // Fetch the resource here; you need to refetch it on every try, since
//         // if you got a conflict on the last update attempt then you need to get
//         // the current version before making your own changes.
//         pod, err := c.Pods("mynamespace").Get(name, metav1.GetOptions{})
//         if err ! nil {
//             return err
//         }
//
//         // Make whatever updates to the resource are needed
//         pod.Status.Phase = v1.PodFailed
//
//         // Try to update
//         _, err = c.Pods("mynamespace").UpdateStatus(pod)
//         // You have to return err itself here (not wrapped inside another error)
//         // so that OnConflict can identify it correctly.
//         return err
//     })
//     if err != nil {
//         // May be conflict if max retries were hit, or may be something unrelated
//         // like permissions or a network error
//         return err
//     }
//     ...
func OnConflict(backoff wait.Backoff, fn func() error) error {
	return OnError(backoff, errors.IsConflict, fn)
}

// OnExists can be used to recreate a resource immediately after the previous copy
// has been deleted, and thus might still be in a state of being removed.
func OnExists(backoff wait.Backoff, fn func() error) error {
	return OnError(backoff, errors.IsAlreadyExists, fn)
}

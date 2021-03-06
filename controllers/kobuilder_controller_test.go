package controllers

import (
	"context"
	"time"

	kov1alpha1 "github.com/feloy/ko-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("KoBuilder controller", func() {

	const timeout = time.Second * 10
	const interval = time.Second * 1

	BeforeEach(func() {
		// Add any setup steps that needs to be executed before each test
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	Context("A KoBuilder resource is created", func() {
		It("ConfigMap and Job should be created successfully", func() {

			key := types.NamespacedName{
				Name:      "my-ko-builder",
				Namespace: "my-ns",
			}

			cmKey := types.NamespacedName{
				Name:      "my-ko-builder-config",
				Namespace: "my-ns",
			}

			jobKey := types.NamespacedName{
				Name:      "my-ko-builder-job",
				Namespace: "my-ns",
			}

			created := &kov1alpha1.KoBuilder{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
				Spec: kov1alpha1.KoBuilderSpec{
					Registry:       "user/ko-builder",
					ServiceAccount: "account@project.com",
					Repository:     "github/com/test/repo",
					Checkout:       "1.2.3",
					ConfigPath:     "/templates",
				},
			}

			// Create
			Expect(k8sClient.Create(context.Background(), created)).Should(Succeed())
			defer k8sClient.Delete(context.Background(), created)

			By("Expecting configmap created with correct data")
			Eventually(func() bool {
				f := &corev1.ConfigMap{}
				return k8sClient.Get(context.Background(), cmKey, f) == nil &&
					f.Data["REGISTRY"] == "user/ko-builder" &&
					f.Data["SERVICE_ACCOUNT"] == "account@project.com" &&
					f.Data["REPOSITORY"] == "github/com/test/repo" &&
					f.Data["CHECKOUT"] == "1.2.3" &&
					f.Data["CONFIG_PATH"] == "/templates"
			}, timeout, interval).Should(BeTrue())

			By("Expecting job created")
			Eventually(func() error {
				f := &batchv1.Job{}
				return k8sClient.Get(context.Background(), jobKey, f)
			}, timeout, interval).Should(BeNil())
		})

		It("KoBuilder status should be Unknown", func() {

			key := types.NamespacedName{
				Name:      "my-ko-builder",
				Namespace: "my-ns",
			}

			created := &kov1alpha1.KoBuilder{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
				Spec: kov1alpha1.KoBuilderSpec{
					Registry:       "user/ko-builder",
					ServiceAccount: "account@project.com",
					Repository:     "github/com/test/repo",
					Checkout:       "1.2.3",
					ConfigPath:     "/templates",
				},
			}

			// Create
			Expect(k8sClient.Create(context.Background(), created)).Should(Succeed())
			defer k8sClient.Delete(context.Background(), created)

			Eventually(func() bool {
				f := &kov1alpha1.KoBuilder{}
				return k8sClient.Get(context.Background(), key, f) == nil &&
					f.Status.State == kov1alpha1.Unknown
			}, timeout, interval).Should(BeTrue())
		})

		Context("The pod of job is active", func() {

			It("KoBuilder status should be Deploying", func() {

				key := types.NamespacedName{
					Name:      "my-ko-builder",
					Namespace: "my-ns",
				}

				jobKey := types.NamespacedName{
					Name:      "my-ko-builder-job",
					Namespace: "my-ns",
				}

				created := &kov1alpha1.KoBuilder{
					ObjectMeta: metav1.ObjectMeta{
						Name:      key.Name,
						Namespace: key.Namespace,
					},
					Spec: kov1alpha1.KoBuilderSpec{
						Registry:       "user/ko-builder",
						ServiceAccount: "account@project.com",
						Repository:     "github/com/test/repo",
						Checkout:       "1.2.3",
						ConfigPath:     "/templates",
					},
				}

				// Create
				Expect(k8sClient.Create(context.Background(), created)).Should(Succeed())
				defer k8sClient.Delete(context.Background(), created)

				By("Expecting job created")
				Eventually(func() error {
					f := &batchv1.Job{}
					return k8sClient.Get(context.Background(), jobKey, f)
				}, timeout, interval).Should(BeNil())

				job := &batchv1.Job{}
				k8sClient.Get(context.Background(), jobKey, job)
				job.Status.Active = 1
				job.Status.Failed = 0
				job.Status.Succeeded = 0
				k8sClient.Status().Update(context.Background(), job)

				Eventually(func() bool {
					f := &kov1alpha1.KoBuilder{}
					return k8sClient.Get(context.Background(), key, f) == nil &&
						f.Status.State == kov1alpha1.Deploying
				}, timeout, interval).Should(BeTrue())
			})
		})

		Context("The pod of job is succeeded", func() {

			It("KoBuilder status should be Deployed", func() {

				key := types.NamespacedName{
					Name:      "my-ko-builder",
					Namespace: "my-ns",
				}

				jobKey := types.NamespacedName{
					Name:      "my-ko-builder-job",
					Namespace: "my-ns",
				}

				created := &kov1alpha1.KoBuilder{
					ObjectMeta: metav1.ObjectMeta{
						Name:      key.Name,
						Namespace: key.Namespace,
					},
					Spec: kov1alpha1.KoBuilderSpec{
						Registry:       "user/ko-builder",
						ServiceAccount: "account@project.com",
						Repository:     "github/com/test/repo",
						Checkout:       "1.2.3",
						ConfigPath:     "/templates",
					},
				}

				// Create
				Expect(k8sClient.Create(context.Background(), created)).Should(Succeed())
				defer k8sClient.Delete(context.Background(), created)

				By("Expecting job created")
				Eventually(func() error {
					f := &batchv1.Job{}
					return k8sClient.Get(context.Background(), jobKey, f)
				}, timeout, interval).Should(BeNil())

				job := &batchv1.Job{}
				k8sClient.Get(context.Background(), jobKey, job)
				job.Status.Succeeded = 1
				job.Status.Failed = 0
				job.Status.Active = 0
				k8sClient.Status().Update(context.Background(), job)

				Eventually(func() bool {
					f := &kov1alpha1.KoBuilder{}
					return k8sClient.Get(context.Background(), key, f) == nil &&
						f.Status.State == kov1alpha1.Deployed
				}, timeout, interval).Should(BeTrue())
			})
		})

		Context("The pod of job is failed", func() {

			It("KoBuilder status should be ErrorDeploying", func() {

				key := types.NamespacedName{
					Name:      "my-ko-builder",
					Namespace: "my-ns",
				}

				jobKey := types.NamespacedName{
					Name:      "my-ko-builder-job",
					Namespace: "my-ns",
				}

				created := &kov1alpha1.KoBuilder{
					ObjectMeta: metav1.ObjectMeta{
						Name:      key.Name,
						Namespace: key.Namespace,
					},
					Spec: kov1alpha1.KoBuilderSpec{
						Registry:       "user/ko-builder",
						ServiceAccount: "account@project.com",
						Repository:     "github/com/test/repo",
						Checkout:       "1.2.3",
						ConfigPath:     "/templates",
					},
				}

				// Create
				Expect(k8sClient.Create(context.Background(), created)).Should(Succeed())
				defer k8sClient.Delete(context.Background(), created)

				By("Expecting job created")
				Eventually(func() error {
					f := &batchv1.Job{}
					return k8sClient.Get(context.Background(), jobKey, f)
				}, timeout, interval).Should(BeNil())

				job := &batchv1.Job{}
				k8sClient.Get(context.Background(), jobKey, job)
				job.Status.Failed = 1
				job.Status.Succeeded = 0
				job.Status.Active = 0
				k8sClient.Status().Update(context.Background(), job)

				Eventually(func() bool {
					f := &kov1alpha1.KoBuilder{}
					return k8sClient.Get(context.Background(), key, f) == nil &&
						f.Status.State == kov1alpha1.ErrorDeploying
				}, timeout, interval).Should(BeTrue())
			})
		})
	})

})

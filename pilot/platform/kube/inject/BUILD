load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "configmap.go",
        "http.go",
        "initializer.go",
        "inject.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//pilot/platform/kube:go_default_library",
        "//pilot/proxy:go_default_library",
        "//pilot/tools/version:go_default_library",
        "@com_github_davecgh_go_spew//spew:go_default_library",
        "@com_github_emicklei_go_restful//:go_default_library",
        "@com_github_ghodss_yaml//:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//ptypes:go_default_library",
        "@com_github_golang_protobuf//ptypes/duration:go_default_library",
        "@io_istio_api//mesh/v1alpha1:go_default_library",
        "@io_istio_api//routing/v1alpha1:go_default_library",
        "@io_k8s_api//apps/v1beta1:go_default_library",
        "@io_k8s_api//batch/v1:go_default_library",
        "@io_k8s_api//batch/v2alpha1:go_default_library",
        "@io_k8s_api//core/v1:go_default_library",
        "@io_k8s_api//extensions/v1beta1:go_default_library",
        "@io_k8s_apimachinery//pkg/api/meta:go_default_library",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:go_default_library",
        "@io_k8s_apimachinery//pkg/fields:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime/schema:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime/serializer:go_default_library",
        "@io_k8s_apimachinery//pkg/types:go_default_library",
        "@io_k8s_apimachinery//pkg/util/strategicpatch:go_default_library",
        "@io_k8s_apimachinery//pkg/util/wait:go_default_library",
        "@io_k8s_apimachinery//pkg/util/yaml:go_default_library",
        "@io_k8s_apimachinery//pkg/watch:go_default_library",
        "@io_k8s_client_go//kubernetes:go_default_library",
        "@io_k8s_client_go//kubernetes/scheme:go_default_library",
        "@io_k8s_client_go//rest:go_default_library",
        "@io_k8s_client_go//tools/cache:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    size = "small",
    srcs = [
        "http_test.go",
        "initializer_test.go",
        "inject_test.go",
    ],
    data = glob(["testdata/*.yaml*"]),
    library = ":go_default_library",
    deps = [
        "//pilot/platform/kube:go_default_library",
        "//pilot/proxy:go_default_library",
        "//pilot/test/util:go_default_library",
        "//pilot/tools/version:go_default_library",
        "@com_github_emicklei_go_restful//:go_default_library",
        "@com_github_ghodss_yaml//:go_default_library",
        "@io_istio_api//routing/v1alpha1:go_default_library",
        "@io_k8s_api//core/v1:go_default_library",
        "@io_k8s_apimachinery//pkg/api/meta:go_default_library",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime:go_default_library",
        "@io_k8s_apimachinery//pkg/runtime/schema:go_default_library",
        "@io_k8s_client_go//kubernetes:go_default_library",
        "@io_k8s_client_go//rest:go_default_library",
    ],
)

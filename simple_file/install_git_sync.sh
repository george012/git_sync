#!/bin/bash

set -e

APP_NAME=git_sync
configFIleName=config.json
APP_DIR="/${APP_NAME}"
# systemd 服务路径
SYSTEMD_PATH="/etc/systemd/system"

# 获取脚本所在目录的绝对路径
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
# 获取脚本所在的上级目录，用于之后的自清理
PARENT_DIR=$(dirname ${SCRIPT_DIR})

echo parent dir is:$PARENT_DIR


# 清理：删除整个上级目录（包括此脚本、相关文件和文件夹）
function clear_install_files() {
    if [ -d "${PARENT_DIR}" ]; then
        # 删除所有以 APP_NAME 开头的文件和文件夹
        find "${PARENT_DIR}" -type f -name "${APP_NAME}*" -exec rm -f {} \;
        find "${PARENT_DIR}" -type d -name "${APP_NAME}*" -empty -delete
        rm -rf ${SCRIPT_DIR}
    fi
}

function install() {

    if [[ -f ${SCRIPT_DIR}/${APP_NAME} ]] && [[ -f ${SCRIPT_DIR}/${APP_NAME}.service ]] && [[ -f ${SCRIPT_DIR}/conf/config.json ]]; then
        # 在根目录下创建服务目录（如果不存在）
        if [ ! -d "${APP_DIR}" ]; then
            mkdir -p ${APP_DIR}
            mkdir -p ${APP_DIR}/logs
            mkdir -p ${APP_DIR}/conf
        fi

        # 安装二进制文件
        TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
        if [[ -f "${APP_DIR}/${APP_NAME}" ]]; then
            mv "${APP_DIR}/${APP_NAME}" "${APP_DIR}/${APP_NAME}.bak_${TIMESTAMP}"
        fi
        cp -f "${SCRIPT_DIR}/${APP_NAME}" "${APP_DIR}"

        # 配置文件
        if [[ -f "${APP_DIR}/conf/${configFIleName}" ]]; then
            mv "${APP_DIR}/conf/${configFIleName}" "${APP_DIR}/conf/${configFIleName}.bak_${TIMESTAMP}"
        fi
        cp -f ${SCRIPT_DIR}/conf/${configFIleName} ${APP_DIR}/conf/${configFIleName}

        # 安装systemd服务文件到 SYSTEMD_PATH 并启动服务
        service_file=${SCRIPT_DIR}/${APP_NAME}.service
        if [[ -f "${service_file}" ]]; then
            cp -f "${service_file}" "${SYSTEMD_PATH}"
            service_name=$(basename "${service_file}" .service)
            systemctl enable "${service_name}"
            systemctl daemon-reload
            systemctl start "${service_name}.service" # 启动服务
        fi
    else
        echo "安装文件不完整"
        exit 1
    fi
}

function updateBinary() {
    read -p "是否更新 ${APP_NAME} [yes/no]：" flag
    if [ -z $flag ]; then
        echo "输入错误" && exit 1
    elif [ "$flag" = "yes" -o "$flag" = "ye" -o "$flag" = "y" ]; then
        # 更新二进制
        TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
        systemctl stop "${APP_NAME}.service"
        if [[ -f "${APP_DIR}/${APP_NAME}" ]]; then
            mv "${APP_DIR}/${APP_NAME}" "${APP_DIR}/${APP_NAME}.bak_${TIMESTAMP}"
        fi
        cp -f "${SCRIPT_DIR}/${APP_NAME}" "${APP_DIR}"

        # 更新systemd服务文件
        service_file=${SCRIPT_DIR}/${APP_NAME}.service
        if [[ -f "${service_file}" ]]; then
            cp -f "${service_file}" "${SYSTEMD_PATH}"
            systemctl daemon-reload
            systemctl restart "${APP_NAME}.service" # 重新启动服务
        fi
    fi
}

function uninstall() {
    read -p "是否卸载 ${APP_NAME} [yes/no]：" flag
    if [ -z "$flag" ]; then
        echo "输入错误" && exit 1
    elif [ "$flag" = "yes" ] || [ "$flag" = "ye" ] || [ "$flag" = "y" ]; then
        for service_file in ${SYSTEMD_PATH}/${APP_NAME}*.service; do
            if [[ -f ${service_file} ]]; then
                service_name=$(basename ${service_file} .service)
                systemctl disable --now ${service_name}
                rm -f ${service_file}
            fi
        done

        rm -rf ${APP_DIR}
        rm -rf /usr/local/lib/libgtgo.so && ldconfig
        systemctl daemon-reload
        echo "卸载 ${APP_NAME} 成功"
    fi
}


echo "============================ ${APP_NAME} ============================"
echo "  1、安装 ${APP_NAME}"
echo "  2、更新 ${APP_NAME}"
echo "  3、卸载 ${APP_NAME}"
echo "  4、更新动态库"
echo "======================================================================"
read -p "$(echo -e "请选择[1-3]：")" choose
case $choose in
1)
    install && wait && clear_install_files
    ;;
2)
    updateBinary && wait && clear_install_files
    ;;
3)
    uninstall && wait && clear_install_files
    ;;
*)
    echo "输入错误，请重新输入！"
    ;;
esac
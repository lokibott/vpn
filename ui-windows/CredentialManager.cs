using Microsoft.Win32.SafeHandles;
using System;
using System.Runtime.InteropServices;

public static class CredentialManager
{
    [DllImport("advapi32.dll", CharSet = CharSet.Auto, SetLastError = true)]
    private static extern bool CredWrite([In] ref Credential userCredential, [In] uint flags);

    [DllImport("advapi32.dll", CharSet = CharSet.Auto, SetLastError = true)]
    private static extern bool CredRead(string target, CredentialType type, int reservedFlag, out IntPtr credentialPtr);

    public static void SaveCredentials(string target, string username, string secret)
    {
        var byteArray = System.Text.Encoding.UTF8.GetBytes(secret);
        var credential = new Credential
        {
            TargetName = target,
            UserName = username,
            CredentialBlob = Marshal.StringToCoTaskMemUni(secret),
            CredentialBlobSize = (uint)secret.Length * 2,
            Type = CredentialType.Generic,
            Persist = CredentialPersist.LocalMachine
        };

        if (!CredWrite(ref credential, 0))
        {
            throw new System.ComponentModel.Win32Exception(Marshal.GetLastWin32Error());
        }
    }
}